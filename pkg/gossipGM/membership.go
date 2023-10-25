package gossipGM

import (
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network/code"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type MembershipManager struct {
	membershipList map[string]*code.MemberInfo
	listMutex      sync.RWMutex
	selfHostName   string
	selfID         string

	IncarnationNumberTrack map[string]int
	forwardRequestBuf      []*code.SuspensionRequest
	mu                     sync.Mutex
}

func NewMembershipManager(selfHostName string) *MembershipManager {
	manager := &MembershipManager{
		membershipList: make(map[string]*code.MemberInfo),
		listMutex:      sync.RWMutex{},
		selfHostName:   selfHostName,

		forwardRequestBuf:      make([]*code.SuspensionRequest, 0),
		IncarnationNumberTrack: make(map[string]int),
	}

	manager.initMembershipList()

	return manager
}

func (m *MembershipManager) initMembershipList() {
	m.updateOrAddMember(m.selfHostName)
}

func (m *MembershipManager) JoinToMembershipList(request *code.JoinRequest) {
	m.updateOrAddMember(request.Host)
}

func (m *MembershipManager) LeaveFromMembershipList(hostname string) {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()

	for k, v := range m.membershipList {
		if strings.Contains(k, hostname) && v.StatusCode == code.Alive {
			v.StatusCode = code.Failed
			m.membershipList[k] = v
			break
		}
	}
}

func (m *MembershipManager) MergeMembershipList(receivedMembershipList map[string]*code.MemberInfo) {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()
	for k, v := range receivedMembershipList {
		if v.StatusCode != code.Alive {
			continue
		}

		if _, ok := m.membershipList[k]; ok {
			if m.membershipList[k].Counter < v.Counter && m.membershipList[k].StatusCode == code.Alive {
				m.membershipList[k] = v
				m.membershipList[k].LocalTime = time.Now()
			}
		} else {
			m.membershipList[k] = v
			m.membershipList[k].LocalTime = time.Now()
		}
	}
}

func (m *MembershipManager) updateOrAddMember(hostname string) {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()

	for k, v := range m.membershipList {
		if v.Hostname == hostname && v.StatusCode == code.Alive {
			delete(m.membershipList, k)
			break
		}
	}

	host := hostname
	timestamp := time.Now()
	uniqueHostID := m.generateUniqueHostID(host, timestamp.Format("20060102150405"))
	m.membershipList[uniqueHostID] = &code.MemberInfo{
		Counter:    1,
		LocalTime:  time.Now(),
		StatusCode: code.Alive,
		Hostname:   hostname,
	}
	if host == m.selfHostName {
		m.selfID = uniqueHostID
	}
}

func (m *MembershipManager) generateUniqueHostID(hostname string, timestamp string) string {
	return fmt.Sprintf("%s-daemon%s", hostname, timestamp)
}

func (m *MembershipManager) IncrementSelfCounter() {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()

	if self, ok := m.membershipList[m.selfID]; ok && self.StatusCode == code.Alive {
		self.Counter += 1
	} else {
		logutil.Logger.Errorf("can not find self member instance or self has been marked as failed")
	}
}

func (m *MembershipManager) GetMembershipList() map[string]*code.MemberInfo {
	m.listMutex.RLock()
	defer m.listMutex.RUnlock()

	copiedList := make(map[string]*code.MemberInfo)
	for k, v := range m.membershipList {
		copiedList[k] = v
	}
	return copiedList
}

func (m *MembershipManager) MarkMembersFailedIfNotUpdated(TFail, TCleanup time.Duration) {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()

	currentTime := time.Now()

	for k, v := range m.membershipList {
		if strings.HasPrefix(k, m.selfHostName) {
			continue
		}
		timeElapsed := currentTime.Sub(v.LocalTime)
		if timeElapsed > TFail && v.StatusCode != code.Failed { // If member is alive or suspected and time elapsed exceeds Tfail
			v.StatusCode = code.Failed // Mark as failed
			logutil.Logger.Infof("Mark member as failed:%s last update time:%s cur time:%v", k, v.LocalTime.String(), time.Now())
			m.membershipList[k] = v
			go m.StartCleanup(k, TCleanup)
		}
	}
}

func (m *MembershipManager) RandomlySelectKNeighbors(k int) []string {
	m.listMutex.RLock()
	defer m.listMutex.RUnlock()

	keys := make([]string, 0, len(m.membershipList))
	for key := range m.membershipList {
		keys = append(keys, key)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	selectedNeighbor := make([]string, 0, k)
	for i := 0; i < k && i < len(keys); i++ {
		key := keys[i]
		member := m.membershipList[key]
		if member.Hostname == m.selfHostName || member.StatusCode == code.Failed {
			continue
		}

		selectedNeighbor = append(selectedNeighbor, member.Hostname)
	}

	return selectedNeighbor
}

func (m *MembershipManager) StartCleanup(targetKey string, TCleanup time.Duration) {
	timer := time.NewTimer(TCleanup)
	<-timer.C
	m.listMutex.Lock()
	delete(m.membershipList, targetKey)
	m.listMutex.Unlock()
	logutil.Logger.Infof("cleanup %s", targetKey)
}
