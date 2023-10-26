package code

import "time"

type RunMode int

const (
	PureGossip RunMode = iota + 1
	GossipWithSuspicion
)

type MethodType int32

const (
	Heartbeat MethodType = iota + 1
	Leave
	ListMember
	ListSelf
	ChangeSuspicion
	Suspicion
)

type MemberStatus int

const (
	Alive MemberStatus = iota + 1
	Failed
	Suspected
)

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
	Counter         int          `json:"counter"`         // Counter for the member
	LocalUpdateTime time.Time    `json:"localUpdateTime"` // Local timestamp
	StatusCode      MemberStatus `json:"status_code"`     // Status code 1(alive), 2(suspect), 3(failed)
	Hostname        string       `json:"hostname"`        // The hostname
}

type HeartbeatRequest struct {
	MemberShipList map[string]*MemberInfo `json:"member_ship_list"`
	Mode           RunMode                `json:"mode"`
	ModeChangeTime int64                  `json:"mode_change_time"`
	SentTimeStamp  int64                  `json:"sent_time_stamp"`
}

type SuspensionInfoType int

const (
	Suspect SuspensionInfoType = iota + 1
	InformAlive
	ConfirmFailed
)

type SuspensionRequest struct {
	TargetID          string             `json:"target_id"`
	InfoType          SuspensionInfoType `json:"info_type"`
	IncarnationNumber int                `json:"incarnation_number"`
}
