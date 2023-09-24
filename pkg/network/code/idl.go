package code

import "time"

type MethodType int32

const (
	Join MethodType = iota + 1
	Heartbeat
	Leave
	ListMember
	ListSelf
	ChangeSuspicion
)

type MemberStatus int

const (
	Alive MemberStatus = iota + 1
	Suspect
	Failed
)

type JoinRequest struct {
	Host string `json:"host"`
}

type ListMemberRequest struct {
	Host string `json:"host"`
}

type ChangeSuspicionRequest struct {
	SuspicionFlag bool      `json:"suspicion_flag"`
	Timestamp     time.Time `json:"time"`
}

type ListSelfRequest struct {
	Host string `json:"host"`
}

type LeaveRequest struct {
	Host string `json:"host"`
}

type MemberInfo struct {
	Counter    int          `json:"counter"`     // Counter for the member
	LocalTime  time.Time    `json:"local_time"`  // Local timestamp
	StatusCode MemberStatus `json:"status_code"` // Status code 1(alive), 2(suspect), 3(failed)
	Hostname   string       `json:"hostname"`    // The hostname
}

type HeartbeatRequest struct {
	MemberShipList map[string]*MemberInfo `json:"sub_member_ship_list"`
	SuspicionFlag  bool                   `json:"suspicion_flag"`
	UpdateTime     time.Time              `json:"update_local_time"`
}
