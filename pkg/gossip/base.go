package gossip

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

	hostname string
	tFail    time.Duration
	tCleanup time.Duration
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
		hostname:          selfHost,
		tFail:             config.GetTFail(),
		tCleanup:          config.GetTCleanup(),
	}
	server.Register(service.Handle)

	return &service
}

func (s *Service) Serve() {
	errChan := s.udpServer.Serve()
	logutil.Logger.Debug("start to receive UDP request!")

	s.joinToGroup()

	//go s.membershipManager.StartSuspicionDetection(time.Second * 2)
	heartbeatTicker := time.NewTicker(config.GetTHeartbeat())
	ticker2 := time.NewTicker(time.Second * 10)
	for {
		select {
		case err := <-errChan:
			panic(err)
		case <-heartbeatTicker.C:
			s.detectionRoutine()
		case <-ticker2.C:
			fmt.Println("----------------------------")
			for k, v := range s.membershipManager.GetMembershipList() {
				value, _ := json.Marshal(v)
				logutil.Logger.Debugf("key:%v member:%v", k, string(value))
			}
			fmt.Println("----------------------------")
		}
	}
}

func (s *Service) HandleJoin(request *code.JoinRequest) {
	logutil.Logger.Debugf("recieve join requet:%v", request.Host)
	s.membershipManager.JoinToMembershipList(request)
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
	}

	return nil
}

func (s *Service) joinToGroup() {
	introducerHost := config.GetIntroducerHost()
	r := code.JoinRequest{Host: s.hostname}
	req := &network.CallRequest{
		MethodName: code.Join,
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
	selectedNeighbors := s.membershipManager.RandomlySelectKMembers(config.GetNumOfGossipPerRound())
	membershipList := s.membershipManager.GetMembershipList()
	go s.membershipManager.MarkMembersFailedIfNotUpdated(s.tFail, s.tCleanup)
	r := code.HeartbeatRequest{MemberShipList: membershipList}
	for _, v := range selectedNeighbors {
		req := network.CallRequest{
			MethodName: code.Heartbeat,
			Request:    r,
			TargetHost: v.Hostname,
		}
		err := s.udpClient.Call(&req)
		if err != nil {
			logutil.Logger.Errorf("send heartbeat failed:%v, host:%v", err, v.Hostname)
		}
	}
}
