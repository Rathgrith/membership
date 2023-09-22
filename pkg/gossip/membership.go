package gossip

import (
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network/code"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type MembershipManager struct {
	membershipList     map[string]*code.MemberInfo
	membershipListLock sync.RWMutex
	suspicionTriggered bool
	suspicionTimeStamp time.Time
}

type Broadcast struct {
	Host         string
	PacketType   string
	BroadcastTTL int
}

type JoinResponse struct {
	Host          string
	PacketType    string
	PacketOutTime time.Time
	PacketData    map[string]*code.MemberInfo
}

func NewMembershipManager() *MembershipManager {
	return &MembershipManager{
		membershipList:     make(map[string]*code.MemberInfo),
		membershipListLock: sync.RWMutex{},
		suspicionTriggered: true,
		suspicionTimeStamp: time.Now(),
	}
}

func (m *MembershipManager) InitMembershipList(hostname string) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	m.UpdateOrAddMember(hostname)
}

func (m *MembershipManager) JoinToMembershipList(request *code.JoinRequest) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	m.UpdateOrAddMember(request.Host)
}

func (m *MembershipManager) OverwriteMembershipList(receivedList map[string]*code.MemberInfo) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	for k := range m.membershipList {
		delete(m.membershipList, k)
	}

	for k, v := range receivedList {
		v.LocalTime = time.Now()
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

func (m *MembershipManager) LeaveFromMembershipList(hostname string) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	for k, v := range m.membershipList {
		if strings.Contains(k, hostname) && v.StatusCode == 1 {
			v.StatusCode = 2
			m.membershipList[k] = v
			break
		}
	}
}

// write me a merger function to merge two membership lists
func (m *MembershipManager) MergeMembershipList(receivedList map[string]*code.MemberInfo) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()
	for k, v := range receivedList {
		if v.StatusCode == 1 {
			if _, ok := m.membershipList[k]; ok {
				if m.membershipList[k].Counter < v.Counter {
					m.membershipList[k] = v
					m.membershipList[k].LocalTime = time.Now()
				}
			} else {
				m.membershipList[k] = v
				m.membershipList[k].LocalTime = time.Now()
			}
		}
	}
}

func (m *MembershipManager) UpdateOrAddMember(hostname string) {
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
		m.membershipList[uniqueHostID] = &code.MemberInfo{
			Counter:    1,
			LocalTime:  time.Now(),
			StatusCode: 1,
			Hostname:   hostname,
		}
	}
}

func (m *MembershipManager) IncrementMembershipCounter() {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	for k, v := range m.membershipList {
		if v.StatusCode == 1 && v.Hostname == getHostname() {
			v.Counter += 1
			m.membershipList[k] = v
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

func (m *MembershipManager) GetMembershipList() map[string]*code.MemberInfo {
	m.membershipListLock.RLock()
	defer m.membershipListLock.RUnlock()

	copiedList := make(map[string]*code.MemberInfo)
	for k, v := range m.membershipList {
		copiedList[k] = v
	}
	return copiedList
}

func (m *MembershipManager) MarkMembersFailedIfNotUpdated(Tfail time.Duration) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	currentTime := time.Now()

	for k, v := range m.membershipList {
		if strings.HasPrefix(k, getHostname()) {
			continue
		}
		timeElapsed := currentTime.Sub(v.LocalTime)
		if timeElapsed > Tfail && v.StatusCode != 2 { // If member is alive or suspected and time elapsed exceeds Tfail
			v.StatusCode = 2 // Mark as failed
			fmt.Println("Marking member as failed:", k)
			m.membershipList[k] = v
		}
	}
}

func (m *MembershipManager) CleanupFailedMembers(Tclean time.Duration) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	currentTime := time.Now()

	for k, v := range m.membershipList {
		// if the k's prefix is the same as the current host, continue
		if strings.HasPrefix(k, getHostname()) {
			continue
		}
		timeElapsed := currentTime.Sub(v.LocalTime)
		if timeElapsed > Tclean && v.StatusCode == 2 { // If member is failed and time elapsed exceeds Tclean
			fmt.Println("Removing failed member:", k)
			delete(m.membershipList, k)
		}
	}
}

func (m *MembershipManager) MarkMembersSuspectedIfNotUpdated(Tsus time.Duration) {
	m.membershipListLock.Lock()
	defer m.membershipListLock.Unlock()

	currentTime := time.Now()

	for k, v := range m.membershipList {
		// if k is currenthostname, return
		if strings.HasPrefix(k, getHostname()) {
			continue
		}
		timeElapsed := currentTime.Sub(v.LocalTime)
		if timeElapsed > Tsus && v.StatusCode == 1 { // If member is alive and time elapsed exceeds Tsus
			v.StatusCode = 3 // Mark as suspected
			m.membershipList[k] = v
			logutil.Logger.Infof("Marking member as suspected: %s", k)
		}
	}
}

func (m *MembershipManager) StartFailureDetection(Tfail time.Duration) {
	ticker := time.NewTicker(Tfail)
	for {
		select {
		case <-ticker.C:
			m.MarkMembersFailedIfNotUpdated(Tfail)
		}
	}
}

func (m *MembershipManager) StartSuspicionDetection(Tsus time.Duration) {
	ticker := time.NewTicker(Tsus)
	for {
		select {
		case <-ticker.C:
			if m.suspicionTriggered {
				m.MarkMembersSuspectedIfNotUpdated(Tsus)
			}
		}
	}
}

func (m *MembershipManager) RandomlySelectKMembers(k int) map[string]code.MemberInfo {
	m.membershipListLock.RLock()
	defer m.membershipListLock.RUnlock()

	if len(m.membershipList) < k {
		selectedMembersMap := make(map[string]code.MemberInfo)
		for key, value := range m.membershipList {
			selectedMembersMap[key] = *value // Dereference the pointer to copy the value
		}

		return selectedMembersMap // or handle the case differently, e.g., return all members in the map
	}

	keys := make([]string, 0, len(m.membershipList))
	for key := range m.membershipList {
		keys = append(keys, key)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	selectedMembersMap := make(map[string]code.MemberInfo)
	for i := 0; i < k; i++ {
		key := keys[i]
		selectedMembersMap[key] = *m.membershipList[key] // Dereference the pointer to copy the value
	}
	// logutil.Logger.Infof("RandomlySelectKMembers: %v", selectedMembersMap)
	return selectedMembersMap
}

func (m *MembershipManager) StartCleanupRoutine(Tcleanup time.Duration) {
	ticker := time.NewTicker(Tcleanup)
	for {
		select {
		case <-ticker.C:
			m.CleanupFailedMembers(Tcleanup)
		}
	}
}

func (m *MembershipManager) EnableSuspicion(requestTime time.Time) {
	if !m.suspicionTriggered {
		if requestTime.After(m.suspicionTimeStamp) {
			m.suspicionTriggered = true
			m.suspicionTimeStamp = requestTime
		}
	}
}

func (m *MembershipManager) DisableSuspicion(requestTime time.Time) {
	if m.suspicionTriggered {
		if requestTime.After(m.suspicionTimeStamp) {
			m.suspicionTriggered = false
			m.suspicionTimeStamp = requestTime
		}
	}
}
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
	// hostname format fa23-cs425-48XX.cs.illinois.edu
}
