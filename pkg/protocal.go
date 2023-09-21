package pkg

import (
	"strings"
	"sync"
	"time"
)

type MembershipManager struct {
	membershipList     map[string]MemberInfo
	membershipListLock sync.RWMutex
}

// Define the structure of member info
type MemberInfo struct {
	Counter    int       // Counter for the member
	LocalTime  time.Time // Local timestamp
	StatusCode int       // Status code 1(alive), 2(suspect), 3(failed)
	Hostname   string    // The hostname
}

type Broadcast struct {
	Host         string
	PacketType   string
	BroadcastTTL int
}

type JoinRequest struct {
	Host          string
	PacketType    string
	PacketOutTime time.Time
	// PacketData    map[int]MemberInfo
}

type JoinResponse struct {
	Host          string
	PacketType    string
	PacketOutTime time.Time
	PacketData    map[string]MemberInfo
}

func NewMembershipManager() *MembershipManager {
	return &MembershipManager{
		membershipList: make(map[string]MemberInfo),
	}
}

func (m *MembershipManager) InitMembershiplist(hostname string) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	m.updateOrAddMember(hostname)
}

func (m *MembershipManager) JoinToMembershipList(request JoinRequest, addr string) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	m.updateOrAddMember(request.Host)
}

func (m *MembershipManager) OverwriteMembershipList(receivedList map[string]MemberInfo) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	for k := range m.membershipList {
		delete(m.membershipList, k)
	}

	for k, v := range receivedList {
		m.membershipList[k] = v
	}
}

func (m *MembershipManager) UpdateLocalTimestampForNode(hostname string) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	for k, v := range m.membershipList {
		if strings.Contains(k, hostname) && v.StatusCode == 1 {
			v.LocalTime = time.Now()
			m.membershipList[k] = v
			break
		}
	}
}

func (m *MembershipManager) updateOrAddMember(hostname string) {
	var activeDaemonKey string

	for k, v := range m.membershipList {
		if v.Hostname == hostname && v.StatusCode == 1 {
			activeDaemonKey = k
			break
		}
	}

	if activeDaemonKey != "" {
		existingMember := m.membershipList[activeDaemonKey]
		existingMember.Counter += 1
		existingMember.LocalTime = time.Now()
		m.membershipList[activeDaemonKey] = existingMember
	} else {
		ipAddr := hostname
		timestamp := time.Now()
		uniqueHostID := ipAddr + "-daemon" + timestamp.Format("20060102150405")
		m.membershipList[uniqueHostID] = MemberInfo{
			Counter:    1,
			LocalTime:  time.Now(),
			StatusCode: 1,
			Hostname:   hostname,
		}
	}
}

func (m *MembershipManager) getPrefixCount(ipPrefix string) int {
	count := 0
	for id := range m.membershipList {
		if strings.HasPrefix(id, ipPrefix) {
			count++
		}
	}
	return count + 1
}

func (m *MembershipManager) GetMembershipList() map[string]MemberInfo {
	m.membershipListLock.RLock()
	defer m.membershipListLock.RUnlock()

	copiedList := make(map[string]MemberInfo)
	for k, v := range m.membershipList {
		copiedList[k] = v
	}
	return copiedList
}
