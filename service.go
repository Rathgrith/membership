package membership

import (
	"encoding/json"
	"fmt"
	"membership/internal/gossipGM"
	"membership/internal/logutil"
	"membership/internal/network"
	"membership/internal/network/code"
	"sync"
	"time"
)

type ServiceHandleFunc func(reqBody []byte) error

type Service struct {
	membershipManager *gossipGM.MembershipManager
	udpServer         *network.CallUDPServer
	udpClient         *network.CallUDPClient
	handleFuncMap     map[code.MethodType]ServiceHandleFunc

	interestFailHost map[string]bool
	interestAll      bool
	failNotifyChan   chan<- string

	hostname            string
	mode                code.RunMode
	modeUpdateTimestamp int64
	runModeMutex        sync.RWMutex
	tHeartbeat          time.Duration
	tFail               time.Duration
	tCleanup            time.Duration
	FanOut              int

	tSuspect time.Duration
	tConfirm time.Duration

	heartbeatCounter int
}

func NewGossipGMService(gmConfig *GossipGMConfig) (*Service, error) {
	selfHost := network.GetSelfHostName()
	manager := gossipGM.NewMembershipManager(selfHost)

	server, err := network.NewUDPServer(gmConfig.ListenPort)
	if err != nil {
		return nil, fmt.Errorf("server boot failed:%w", err)
	}
	client := network.NewCallUDPClient()

	service := Service{
		membershipManager:   manager,
		udpServer:           server,
		udpClient:           client,
		hostname:            selfHost,
		tHeartbeat:          gmConfig.THeartbeat,
		tFail:               gmConfig.TFail,
		tCleanup:            gmConfig.TCleanup,
		FanOut:              gmConfig.OutPerRound,
		mode:                gmConfig.Mode,
		modeUpdateTimestamp: 0, // wait running member's heartbeat to update
		runModeMutex:        sync.RWMutex{},
		interestAll:         false,
		interestFailHost:    map[string]bool{},
	}
	service.initHandleFunc()
	server.Register(service.handle)

	return &service, nil
}

func (s *Service) handleHeartbeat(reqBody []byte) error {
	req := code.HeartbeatRequest{}
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		return err
	}

	// TODO: provide separate run mode update interface
	s.runModeMutex.RLock()
	if req.ModeChangeTime > s.modeUpdateTimestamp { // update run mode
		s.runModeMutex.Unlock()
		s.updateRunMode(req.Mode, req.ModeChangeTime)
	} else {
		s.runModeMutex.RUnlock()
	}

	s.membershipManager.MergeMembershipList(req.MemberShipList)

	return nil
}

func (s *Service) updateRunMode(newMode code.RunMode, updateTimestamp int64) {
	if updateTimestamp <= s.modeUpdateTimestamp { // check whether the update condition still valid
		return
	}
	s.runModeMutex.Lock()
	s.mode = newMode
	s.modeUpdateTimestamp = updateTimestamp
	s.runModeMutex.Unlock()
}

func (s *Service) handleLeave(reqBody []byte) error {
	// TODO: handle leave
	return nil
}

func (s *Service) handleListMember(reqBody []byte) error {
	logutil.Logger.Infof("Listing all the Members........................")
	for k, v := range s.membershipManager.GetMembershipList() {
		logutil.Logger.Infof("member ID: %v, Attributes: %v", k, v)
	}
	return nil
}

func (s *Service) handleListSelf(reqBody []byte) error {
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
		code.Heartbeat:  s.handleHeartbeat,
		code.ListMember: s.handleListMember,
		code.ListSelf:   s.handleListSelf,
	}
}

func (s *Service) handle(header *code.RequestHeader, reqBody []byte) error {
	f, ok := s.handleFuncMap[header.Method]
	if !ok {
		return fmt.Errorf("unknown Method:%v", header.Method)
	}
	return f(reqBody)
}

func (s *Service) heartbeat(membershipList map[string]*code.MemberInfo,
	hostsOfTargets []string, piggybackRequests []*network.CallRequest) {
	heartBeat := &code.HeartbeatRequest{
		MemberShipList: membershipList,
		Mode:           s.mode,
		ModeChangeTime: s.modeUpdateTimestamp,
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
	selectedNeighbors := s.membershipManager.RandomlySelectKNeighborsHost(s.FanOut)
	membershipList := s.membershipManager.GetMembershipList()
	go func() {
		failedMemberHost := s.membershipManager.MarkMembersFailedIfNotUpdated(s.tFail, s.tCleanup)
		for _, host := range failedMemberHost {
			if s.interestAll || s.interestFailHost[host] {
				s.failNotifyChan <- host
			}
		}
	}()
	s.heartbeat(membershipList, selectedNeighbors, nil)
}

func (s *Service) routine() {
	s.pureGossipRoutine()
	s.heartbeatCounter++
}
