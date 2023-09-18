package pkg

import "time"

// Define the structure of member info
type MemberInfo struct {
	Counter    int       // Counter for the member
	LocalTime  time.Time // Local timestamp
	StatusCode int       // Status code (you can specify the type further if needed)
}

type JoinRequest struct {
	HostID        int
	PacketType    string
	PacketOutTime time.Time
	// PacketData    map[int]MemberInfo
}

type JoinResponse struct {
	HostID            int
	PacketType        string
	PacketDataOutTime time.Time
	PacketData        map[int]MemberInfo
}
