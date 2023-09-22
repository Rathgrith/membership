package gossip

import (
	"ece428_mp2/pkg"
	"fmt"
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

}

func (s *Service) Test() {
	fmt.Println(s.membershipManager.GetMembershipList())
}
