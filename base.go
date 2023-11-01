package membership

import (
	"membership/internal/logutil"
	"membership/internal/network/code"
	"time"
)

func NewDefaultGossipGMService() (*Service, error) {
	return NewGossipGMService(GetDefaultGossipGMConfig())
}

func (s *Service) Serve() {
	errChan := s.udpServer.Serve()
	logutil.Logger.Debug("start to receive UDP request!")

	heartbeatTicker := time.NewTicker(s.tHeartbeat)
	for {
		select {
		case err := <-errChan:
			logutil.Logger.Errorf(err.Error())
		case <-heartbeatTicker.C:
			s.routine()
		}
	}
}

func (s *Service) JoinToGroup(introducerHostList []string) {
	s.heartbeat(s.membershipManager.GetMembershipList(), introducerHostList, nil)
}

func (s *Service) GetHostsOfAllMembers() []string {
	list := s.membershipManager.GetMembershipList()
	hosts := make([]string, 0, len(list))
	for _, v := range list {
		if v.StatusCode == code.Failed {
			continue
		}
		hosts = append(hosts, v.Hostname)
	}
	return hosts
}

func (s *Service) SubscribeFailNotification(interestHost []string, all bool, notifyChan chan<- string) {
	if all {
		s.interestAll = true
	} else {
		m := make(map[string]bool, len(interestHost))
		for _, host := range interestHost {
			m[host] = true
		}
		s.interestFailHost = m
	}

	s.failNotifyChan = notifyChan
}
