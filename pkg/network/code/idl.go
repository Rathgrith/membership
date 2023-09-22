package code

import "time"

type MethodType int32

const (
	Join MethodType = iota + 1
	Heartbeat
	Leave
	ListMember
	ListSelf
)

type JoinRequest struct {
	Host string `json:"host"`
}

type List_Member struct {
	Host string `json:"host"`
}

type List_Self struct {
	Host string `json:"host"`
}

type Leave_Request struct {
	Host string `json:"host"`
}

// Define the structure of member info
type MemberInfo struct {
	Counter    int       `json:"counter"`     // Counter for the member
	LocalTime  time.Time `json:"local_time"`  // Local timestamp
	StatusCode int       `json:"status_code"` // Status code 1(alive), 2(suspect), 3(failed)
	Hostname   string    `json:"hostname"`    // The hostname
}

type HeartbeatRequest struct {
	MemberShipList map[string]*MemberInfo `json:"sub_member_ship_list"`
	SuspicionFlag  bool                   `json:"suspicion_flag"`
	UpdateTime     time.Time              `json:"update_local_time"`
}
