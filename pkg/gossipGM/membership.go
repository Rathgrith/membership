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
}

func NewMembershipManager(selfHostName string) *MembershipManager {
	manager := &MembershipManager{
		membershipList: make(map[string]*code.MemberInfo),
		listMutex:      sync.RWMutex{},
		selfHostName:   selfHostName,
	}

	manager.initMembershipList()

	return manager
}

func (m *MembershipManager) initMembershipList() {
	m.addSelfToList(m.selfHostName)
}

func (m *MembershipManager) addSelfToList(hostname string) {
	// this function will only be called when init, do not use mutex to protect write!
	host := hostname
	timestamp := time.Now()
	uniqueHostID := m.generateUniqueHostID(host, timestamp.Format("20060102150405"))
	m.membershipList[uniqueHostID] = &code.MemberInfo{
		Counter:         1,
		LocalUpdateTime: time.Now(),
		StatusCode:      code.Alive,
		Hostname:        hostname,
	}
	m.selfID = uniqueHostID
	logutil.Logger.Infof("self id:%v", uniqueHostID)
}

func (m *MembershipManager) LeaveFromMembershipList(hostname string) {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()

	for k, v := range m.membershipList {
		if strings.Contains(k, hostname) && v.StatusCode == code.Alive {
			m.membershipList[k].StatusCode = code.Failed
			go m.StartCleanup(k, 10)
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

		if _, ok := m.membershipList[k]; ok { // update
			if m.membershipList[k].Counter < v.Counter && m.membershipList[k].StatusCode == code.Alive {
				m.membershipList[k].Counter = v.Counter
				m.membershipList[k].LocalUpdateTime = time.Now()
			}
		} else { // add
			logutil.Logger.Debugf("add %v to membership list", k)
			m.membershipList[k] = v
			m.membershipList[k].LocalUpdateTime = time.Now()
		}
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
		self.LocalUpdateTime = time.Now()
	} else {
		logutil.Logger.Errorf("can not find self member instance or self has been marked as failed")
	}
}

func (m *MembershipManager) GetMembershipList() map[string]*code.MemberInfo {
	m.listMutex.RLock()
	defer m.listMutex.RUnlock()

	copiedList := make(map[string]*code.MemberInfo)
	for k, v := range m.membershipList { // shallow copy
		copiedList[k] = v
	}
	return copiedList
}

func (m *MembershipManager) MarkMembersFailedIfNotUpdated(TFail, TCleanup time.Duration) {
	m.listMutex.Lock()
	defer m.listMutex.Unlock()

	currentTime := time.Now()

	for k, v := range m.membershipList {
		timeElapsed := currentTime.Sub(v.LocalUpdateTime)
		if timeElapsed > TFail && v.StatusCode != code.Failed {
			m.membershipList[k].StatusCode = code.Failed
			logutil.Logger.Infof("Mark member:%s as failed, last update time:%s, elapsed:%v",
				k, v.LocalUpdateTime.String(), time.Now().Sub(v.LocalUpdateTime))
			go m.StartCleanup(k, TCleanup)
		}
	}
}

func (m *MembershipManager) RandomlySelectKNeighbors(k int) []string {
	keys := make([]string, 0, len(m.membershipList))
	selectedNeighbor := make([]string, 0, k)

	m.listMutex.RLock()
	defer m.listMutex.RUnlock()
	for key := range m.membershipList {
		keys = append(keys, key)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

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
