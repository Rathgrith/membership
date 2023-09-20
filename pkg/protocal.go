package pkg

import (
	// "ece428_mp2/config"
	"fmt"
	"net"
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
}

type JoinRequest struct {
	HostID        string
	PacketType    string
	PacketOutTime time.Time
	// PacketData    map[int]MemberInfo
}

type JoinResponse struct {
	HostID        string
	PacketType    string
	PacketOutTime time.Time
	PacketData    map[string]MemberInfo
}

func InitMembershiplist(hostname string) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()
	ipAddr := hostname
	prefixCount := getPrefixCount(ipAddr)

	// Create a unique HostID using IP address and prefix count
	uniqueHostID := ipAddr + "-daemon" + fmt.Sprintf("%d", prefixCount)

	membershipList[uniqueHostID] = MemberInfo{
		Counter:    1,
		LocalTime:  time.Now(),
		StatusCode: 1,
	}
}

func JoinToMembershipList(request JoinRequest, addr net.Addr) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()

	// Extract IP address without port
	ipAddr := strings.Split(addr.String(), ":")[0]
	prefixCount := getPrefixCount(ipAddr)

	// Create a unique HostID using IP address and prefix count
	uniqueHostID := ipAddr + "-daemon" + fmt.Sprintf("%d", prefixCount)

	membershipList[uniqueHostID] = MemberInfo{
		Counter:    1,
		LocalTime:  time.Now(),
		StatusCode: 1,
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

func UpdateMembershipList(receivedList map[string]MemberInfo) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()

	for k, v := range receivedList {
		// If the key exists in our current membershipList
		if existingMember, ok := membershipList[k]; ok {
			// Choose the larger counter between the existing member and the received member,
			// then increment it by 1.
			if existingMember.Counter > v.Counter {
				v.Counter = existingMember.Counter + 1
			} else {
				v.Counter += 1
			}

			// Update the member's timestamp to the newer one
			if v.LocalTime.After(existingMember.LocalTime) {
				existingMember.LocalTime = v.LocalTime
			}

			membershipList[k] = v
		} else {
			// Otherwise, add the received member info to our list
			membershipList[k] = v
		}
	}
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
