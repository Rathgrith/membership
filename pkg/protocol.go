package pkg

import (
	"ece428_mp2/config"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	membershipList = make(map[int]MemberInfo)
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
	HostID        int
	PacketType    string
	PacketOutTime time.Time
	// PacketData    map[int]MemberInfo
}

type JoinResponse struct {
	HostID        int
	PacketType    string
	PacketOutTime time.Time
	PacketData    map[int]MemberInfo
}

func InitMembershiplist(hostname string) {
	// Initialize membership list
	// selfHost = "fa23-cs425-4810.cs.illinois.edu"
	selfID, err := config.GetHostID(hostname)
	if err != nil {
		fmt.Println("Error:", err)
	}
	membershipList[selfID] = MemberInfo{
		Counter:    1,
		LocalTime:  time.Now(),
		StatusCode: 1,
	}
}

func JoinToMembershipList(request JoinRequest, addr net.Addr) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()
	membershipList[request.HostID] = MemberInfo{
		Counter:    1,
		LocalTime:  time.Now(),
		StatusCode: 1,
	}
}

func UpdateMembershipList(receivedList map[int]MemberInfo) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()

	for k, v := range receivedList {
		if existingMember, ok := membershipList[k]; ok && v.StatusCode == 1 {
			existingMember.Counter += v.Counter
			existingMember.LocalTime = v.LocalTime
			membershipList[k] = existingMember
		} else {
			membershipList[k] = v
		}
	}
}

func GetMembershipList() map[int]MemberInfo {
	membershipListLock.RLock()
	defer membershipListLock.RUnlock()

	copiedList := make(map[int]MemberInfo)
	for k, v := range membershipList {
		copiedList[k] = v
	}
	return copiedList
}
