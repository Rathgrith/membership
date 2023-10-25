package gossipGM

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	membershipManager *MembershipManager
	udpServer         *network.CallUDPServer
	udpClient         *network.CallUDPClient

	hostname  string
	mode      code.RunMode
	timeStamp time.Time
	tFail     time.Duration
	tCleanup  time.Duration
	tSuspect  time.Duration
	tConfirm  time.Duration

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
		membershipManager: manager,
		udpServer:         server,
		udpClient:         client,
		timeStamp:         time.Time{},
		hostname:          selfHost,
		tFail:             config.GetTFail(),
		tCleanup:          config.GetTCleanup(),
		mode:              config.GetDefaultRunMode(),
		tSuspect:          config.GetTSuspect(),
		tConfirm:          config.GetTConfirm(),
	}
	server.Register(service.Handle)

	return &service
}

func (s *Service) Serve() {
	errChan := s.udpServer.Serve()
	logutil.Logger.Debug("start to receive UDP request!")

	s.joinToGroup()

	heartbeatTicker := time.NewTicker(config.GetTHeartbeat())
	for {
		select {
		case err := <-errChan:
			logutil.Logger.Errorf(err.Error())
		case <-heartbeatTicker.C:
			s.detectionRoutine()
		}
	}
}

func (s *Service) HandleRunModeChange(flag bool, timestamp time.Time) {
	if timestamp.After(s.timeStamp) {
		if flag == false {
			s.mode = code.PureGossip
		} else {
			s.mode = code.GossipWithSuspicion
		}
	}
	fmt.Println("suspicion flag changed to:", s.mode)
}

func (s *Service) HandleJoin(request *code.JoinRequest) {
	logutil.Logger.Infof("receive join request:%v", request.Host)
	s.membershipManager.JoinToMembershipList(request)
}

func (s *Service) HandleSuspicion(request *code.SuspensionRequest) {
	s.membershipManager.HandleSuspicionRequest(request)
}

func (s *Service) HandleLeave() {
	s.membershipManager.LeaveFromMembershipList(config.GetSelfHostName())
}

func (s *Service) ListMember() {
	logutil.Logger.Infof("Listing all the Members........................")
	for k, v := range s.membershipManager.GetMembershipList() {
		logutil.Logger.Infof("member ID: %v, Attributes: %v", k, v)
	}
}

func (s *Service) ListSelf(hostname string) {
	// match the member with the same hostname
	logutil.Logger.Infof("Listing self........................")
	for k, v := range s.membershipManager.GetMembershipList() {
		if v.Hostname == hostname {
			logutil.Logger.Infof("//////////current member ID: %v", k)
			logutil.Logger.Infof("current heartbeat counter: %v", s.heartbeatCounter)
		}
	}
}

func (s *Service) Handle(header *code.RequestHeader, reqBody []byte) error {
	if header.Method == code.Join {
		req := code.JoinRequest{}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			return err
		}
		s.HandleJoin(&req)
	} else if header.Method == code.Heartbeat {
		req := code.HeartbeatRequest{}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			return err
		}
		if req.UpdateTime.After(s.timeStamp) {
			logutil.Logger.Infof("update timestamp:%v", req.UpdateTime)
			logutil.Logger.Infof("update mode:%v", req.SuspicionFlag)
			if req.SuspicionFlag == false {
				s.mode = code.PureGossip
			} else {
				s.mode = code.GossipWithSuspicion
			}
			s.timeStamp = req.UpdateTime
		}
		s.membershipManager.MergeMembershipList(req.MemberShipList)
	} else if header.Method == code.ListMember {
		s.ListMember()
	} else if header.Method == code.Leave {
		s.HandleLeave()
	} else if header.Method == code.ListSelf {
		req := code.HeartbeatRequest{}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			return err
		}
		s.ListSelf(s.hostname)
	} else if header.Method == code.ChangeSuspicion {
		req := code.ChangeSuspicionRequest{}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			return err
		}
		if req.Timestamp.After(s.timeStamp) {
			s.timeStamp = req.Timestamp
			if req.SuspicionFlag == false {
				s.mode = code.PureGossip
			} else {
				s.mode = code.GossipWithSuspicion
			}
		}
		s.HandleRunModeChange(req.SuspicionFlag, req.Timestamp)
	} else if header.Method == code.Suspicion {
		req := code.SuspensionRequest{}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			return err
		}
		s.HandleSuspicion(&req)
	}

	return nil
}

func (s *Service) joinToGroup() {
	introducerHost := config.GetIntroducerHost()
	list := s.membershipManager.GetMembershipList()
	r := code.HeartbeatRequest{
		MemberShipList: list,
		SuspicionFlag:  false,
		UpdateTime:     time.Time{},
	}
	req := &network.CallRequest{
		MethodName: code.Heartbeat,
		Request:    r,
		TargetHost: introducerHost,
	}
	err := s.udpClient.Call(req)
	if err != nil {
		panic(err)
	}
	logutil.Logger.Debug("join request sent")
}

func (s *Service) detectionRoutine() {
	s.membershipManager.IncrementSelfCounter()
	selectedNeighbors := s.membershipManager.RandomlySelectKNeighbors(config.GetNumOfGossipPerRound())
	var forwardRequests []*code.SuspensionRequest
	if s.mode == code.GossipWithSuspicion {
		s.membershipManager.MarkMembersSuspectedIfNotUpdated(s.tSuspect, s.tConfirm)
		forwardRequests = s.membershipManager.GetAllForwardSuspicionRequest()
	}
	membershipList := s.membershipManager.GetMembershipList()
	flag := false
	if s.mode == code.PureGossip {
		go s.membershipManager.MarkMembersFailedIfNotUpdated(s.tFail, s.tCleanup)
	} else {
		flag = true
	}
	r := code.HeartbeatRequest{
		MemberShipList: membershipList,
		UpdateTime:     s.timeStamp,
		SuspicionFlag:  flag,
	}
	for _, neighborHost := range selectedNeighbors {
		for _, fr := range forwardRequests {
			req := network.CallRequest{
				MethodName: code.Suspicion,
				Request:    fr,
				TargetHost: neighborHost,
			}
			err := s.udpClient.Call(&req)
			if err != nil {
				logutil.Logger.Errorf("forward Suspicion failed:%v, host:%v", err, neighborHost)
			}
		}
		req := network.CallRequest{
			MethodName: code.Heartbeat,
			Request:    r,
			TargetHost: neighborHost,
		}
		err := s.udpClient.Call(&req)
		if err != nil {
			logutil.Logger.Errorf("send heartbeat failed:%v, host:%v", err, neighborHost)
		}
	}
	s.heartbeatCounter++
}
