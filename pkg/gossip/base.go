package gossip

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"
	"encoding/json"
	"time"
)

type Service struct {
	membershipManager *MembershipManager
	udpServer         *network.CallUDPServer
	udpClient         *network.CallUDPClient
}

func NewGossipService() *Service {
	manager := NewMembershipManager()
	manager.InitMembershipList(getHostname())
	logutil.Logger.Debugf("member:%v", manager.membershipList)

	server, err := network.NewUDPServer(10088) //TODO: read port from config yaml
	if err != nil {
		logutil.Logger.Errorf("server boot failed:%v", err)
		panic(err)
	}
	client := network.NewCallUDPClient()

	service := Service{
		membershipManager: manager,
		udpServer:         server,
		udpClient:         client,
	}
	server.Register(service.Handle)

	return &service
}

func (s *Service) Serve() {
	errChan := s.udpServer.Serve()
	logutil.Logger.Debug("start to receive UDP request!")

	s.joinToGroup()

	go s.membershipManager.StartFailureDetection(time.Second * 3) // Assuming Tfail is 2 seconds
	go s.membershipManager.StartCleanupRoutine(time.Second * 5)   // Assuming Tcleanup is 4 seconds
	//go s.membershipManager.StartSuspicionDetection(time.Second * 2)
	heartbeatTicker := time.NewTicker(config.GetTHeartbeat())
	ticker2 := time.NewTicker(time.Second * 5)
	for {
		select {
		case err := <-errChan:
			panic(err)
		case <-heartbeatTicker.C:
			s.detectionRoutine()
		case <-ticker2.C:
			logutil.Logger.Debugf("----------------------------")
			for k, v := range s.membershipManager.GetMembershipList() {
				value, _ := json.Marshal(v)
				logutil.Logger.Debugf("key:%v member:%v", k, string(value))
			}
			logutil.Logger.Debugf("----------------------------")
		}
	}
}

func (s *Service) Join(request *code.JoinRequest) {
	logutil.Logger.Debugf("recieve join requet:%v", request.Host)
	s.membershipManager.JoinToMembershipList(request)
}

func (s *Service) Leave() {
	s.membershipManager.LeaveFromMembershipList(getHostname())
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
		s.Join(&req)
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
		s.Leave()
	} else if header.Method == code.ListSelf {
		req := code.HeartbeatRequest{}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			return err
		}
		s.ListSelf(getHostname())
	}

	return nil
}

func (s *Service) joinToGroup() {
	introducerHost := config.GetIntroducerHost()
	r := code.JoinRequest{Host: getHostname()}
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
	s.membershipManager.IncrementMembershipCounter()
	list := s.membershipManager.GetMembershipList()
	l := s.membershipManager.RandomlySelectKMembers(3)
	r := code.HeartbeatRequest{MemberShipList: list}
	for _, v := range l {
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
