package gossipGM

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type ServiceHandleFunc func(reqBody []byte) error

type Service struct {
	membershipManager *MembershipManager
	udpServer         *network.CallUDPServer
	udpClient         *network.CallUDPClient
	handleFuncMap     map[code.MethodType]ServiceHandleFunc

	hostname            string
	mode                code.RunMode
	modeUpdateTimestamp int64
	runModeMutex        sync.RWMutex
	tFail               time.Duration
	tCleanup            time.Duration
	tSuspect            time.Duration
	tConfirm            time.Duration

	heartbeatCounter int
}

func NewGossipService() *Service {
	selfHost := config.GetSelfHostName()
	manager := NewMembershipManager(selfHost)

	server, err := network.NewUDPServer(config.GetListenPort())
	if err != nil {
		logutil.Logger.Errorf("server boot failed:%v", err)
		panic(err)
	}
	client := network.NewCallUDPClient()

	service := Service{
		membershipManager:   manager,
		udpServer:           server,
		udpClient:           client,
		hostname:            selfHost,
		tFail:               config.GetTFail(),
		tCleanup:            config.GetTCleanup(),
		mode:                config.GetDefaultRunMode(),
		modeUpdateTimestamp: 0, // wait running member's heartbeat to update
		runModeMutex:        sync.RWMutex{},
		tSuspect:            config.GetTSuspect(),
		tConfirm:            config.GetTConfirm(),
	}
	service.initHandleFunc()
	server.Register(service.Handle)
	//network.CleanUDPReceiveBuffer()

	return &service
}

func (s *Service) Serve() {
	errChan := s.udpServer.Serve()
	logutil.Logger.Debug("start to receive UDP request!")

	s.HandleListMember(nil)
	s.joinToGroup()

	heartbeatTicker := time.NewTicker(config.GetTHeartbeat())
	for {
		select {
		case err := <-errChan:
			logutil.Logger.Errorf(err.Error())
		case <-heartbeatTicker.C:
			s.routine()
		}
	}
}

func (s *Service) HandleHeartbeat(reqBody []byte) error {
	req := code.HeartbeatRequest{}
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		return err
	}

	// TODO: provide separate run mode update interface
	s.runModeMutex.RLock()
	if req.ModeChangeTime > s.modeUpdateTimestamp { // update run mode
		s.runModeMutex.Unlock()
		s.UpdateRunMode(req.Mode, req.ModeChangeTime)
	} else {
		s.runModeMutex.RUnlock()
	}

	logutil.Logger.Debugf("received membership list:%v, sent time:%v", req.MemberShipList, req.SentTimeStamp)
	s.membershipManager.MergeMembershipList(req.MemberShipList)

	return nil
}

func (s *Service) UpdateRunMode(newMode code.RunMode, updateTimestamp int64) {
	if updateTimestamp <= s.modeUpdateTimestamp { // check whether the update condition still valid
		return
	}
	s.runModeMutex.Lock()
	s.mode = newMode
	s.modeUpdateTimestamp = updateTimestamp
	s.runModeMutex.Unlock()
}

func (s *Service) HandleLeave(reqBody []byte) error {
	// TODO: handle leave
	return nil
}

func (s *Service) HandleListMember(reqBody []byte) error {
	logutil.Logger.Infof("Listing all the Members........................")
	for k, v := range s.membershipManager.GetMembershipList() {
		logutil.Logger.Infof("member ID: %v, Attributes: %v", k, v)
	}
	return nil
}

func (s *Service) HandleListSelf(reqBody []byte) error {
	// match the member with the same hostname
	logutil.Logger.Infof("Listing self........................")
	for k, v := range s.membershipManager.GetMembershipList() {
		if v.Hostname == s.hostname {
			logutil.Logger.Infof("current member ID: %v", k)
			logutil.Logger.Infof("current heartbeat counter: %v", s.heartbeatCounter)
		}
	}
	return nil
}

func (s *Service) initHandleFunc() {
	s.handleFuncMap = map[code.MethodType]ServiceHandleFunc{
		code.Heartbeat:  s.HandleHeartbeat,
		code.ListMember: s.HandleListMember,
		code.ListSelf:   s.HandleListSelf,
	}
}

func (s *Service) Handle(header *code.RequestHeader, reqBody []byte) error {
	f, ok := s.handleFuncMap[header.Method]
	if !ok {
		return fmt.Errorf("unknown Method:%v", header.Method)
	}
	return f(reqBody)
}

func (s *Service) joinToGroup() {
	groupIntroducerHost := config.GetIntroducerHost()
	list := s.membershipManager.GetMembershipList()
	r := code.HeartbeatRequest{
		MemberShipList: list,
	}
	req := &network.CallRequest{
		MethodName: code.Heartbeat,
		Request:    r,
		TargetHost: groupIntroducerHost,
	}
	err := s.udpClient.Call(req)
	if err != nil {
		panic(err)
	}
	logutil.Logger.Debugf("join request sent to introducer:%s", groupIntroducerHost)
}

func (s *Service) heartbeat(membershipList map[string]*code.MemberInfo,
	hostsOfTargets []string, piggybackRequests []*network.CallRequest) {
	heartBeat := code.HeartbeatRequest{
		MemberShipList: membershipList,
		Mode:           s.mode,
		ModeChangeTime: s.modeUpdateTimestamp,
		SentTimeStamp:  time.Now().Unix(),
	}

	for _, neighborHost := range hostsOfTargets {
		for _, r := range piggybackRequests {
			err := s.udpClient.Call(r)
			if err != nil {
				logutil.Logger.Errorf("piggyback request failed:%v, req:%v, target host:%v", err, r.MethodName, neighborHost)
			}
		}
		req := network.CallRequest{
			MethodName: code.Heartbeat,
			Request:    heartBeat,
			TargetHost: neighborHost,
		}
		err := s.udpClient.Call(&req)
		if err != nil {
			logutil.Logger.Errorf("send heartbeat failed:%v, host:%v", err, neighborHost)
		}
	}
}

func (s *Service) pureGossipRoutine() {
	s.membershipManager.IncrementSelfCounter()
	selectedNeighbors := s.membershipManager.RandomlySelectKNeighbors(config.GetNumOfGossipPerRound())
	membershipList := s.membershipManager.GetMembershipList()
	go s.membershipManager.MarkMembersFailedIfNotUpdated(s.tFail, s.tCleanup)
	s.heartbeat(membershipList, selectedNeighbors, nil)
}

func (s *Service) routine() {
	s.pureGossipRoutine()
	s.heartbeatCounter++
}
