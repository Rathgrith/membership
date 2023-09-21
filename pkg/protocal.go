package pkg

import (
	// "ece428_mp2/config"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	membershipList = make(map[string]MemberInfo)
	// read/write lock for membership list
	membershipListLock sync.RWMutex
)

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

func InitMembershiplist(hostname string) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()

	updateOrAddMember(hostname)
}

func JoinToMembershipList(request JoinRequest, addr string) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()

	updateOrAddMember(request.Host)
}

func OverwriteMembershipList(receivedList map[string]MemberInfo) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()

	// Clear the current membership list
	for k := range membershipList {
		delete(membershipList, k)
	}

	// Populate the current membership list with the received list
	for k, v := range receivedList {
		membershipList[k] = v
	}
}

func UpdateLocalTimestampForNode(hostname string) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()

	for k, v := range membershipList {
		// Check if the key contains the node's hostname and the node is alive
		if strings.Contains(k, hostname) && v.StatusCode == 1 {
			// Update the LocalTime to the current time for the node
			v.LocalTime = time.Now()
			membershipList[k] = v
			break
		}
	}
}

func updateOrAddMember(hostname string) {
	// Check if a member with the same hostname exists
	var existingDaemonKey string

	for k, v := range membershipList {
		if v.Hostname == hostname {
			existingDaemonKey = k
			break
		}
	}

	if existingDaemonKey != "" && membershipList[existingDaemonKey].StatusCode == 1 {
		// Update the existing daemon if it's active
		existingMember := membershipList[existingDaemonKey]
		existingMember.Counter += 1
		existingMember.LocalTime = time.Now()
		existingMember.StatusCode = 1
		membershipList[existingDaemonKey] = existingMember
	} else {
		// Create a new daemon
		ipAddr := hostname
		prefixCount := getPrefixCount(ipAddr)
		uniqueHostID := ipAddr + "-daemon" + fmt.Sprintf("%d", prefixCount)

		membershipList[uniqueHostID] = MemberInfo{
			Counter:    1,
			LocalTime:  time.Now(),
			StatusCode: 1,
			Hostname:   hostname,
		}
	}
}

func getPrefixCount(ipPrefix string) int {
	count := 0
	for id := range membershipList {
		if strings.HasPrefix(id, ipPrefix) {
			count++
		}
	}
	return count + 1
}

func GetMembershipList() map[string]MemberInfo {
	membershipListLock.RLock()
	defer membershipListLock.RUnlock()

	copiedList := make(map[string]MemberInfo)
	for k, v := range membershipList {
		copiedList[k] = v
	}
	return copiedList
}
