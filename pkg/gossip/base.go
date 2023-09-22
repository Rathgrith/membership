package gossip

import (
	"ece428_mp2/pkg"
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Service struct {
	membershipManager *pkg.MembershipManager
}

func NewGossipService(membershipManager *pkg.MembershipManager) *Service {
	fmt.Println("NewGossipService")
	fmt.Println(membershipManager.GetMembershipList())
	return &Service{
		membershipManager: membershipManager,
	}
}

func (s *Service) Serve() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	server, err := network.NewUDPServer(10088)
	if err != nil {
		logutil.Logger.Error(err)
		panic(err)
	}
	server.Register(s.Handle)

	errChan := server.Serve()

	logutil.Logger.Debug("server started!")
	select {
	case err = <-errChan:
		panic(err)
	}
}

func (s *Service) Join(request *code.JoinRequest) {
	s.membershipManager.JoinToMembershipList(request)
	logutil.Logger.Println(s.membershipManager.GetMembershipList())
}

func (s *Service) Test() {
	fmt.Println(s.membershipManager.GetMembershipList())
}

func (s *Service) Handle(header *code.RequestHeader, reqBody []byte) error {
	if header.Method == code.Join {
		req := code.JoinRequest{}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			return err
		}
		s.Join(&req)
	}

	return nil
}
