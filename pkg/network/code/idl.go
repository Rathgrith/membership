package code

type MethodType int32

const (
	Join MethodType = iota + 1
	Heartbeat
	Leave
)

type JoinRequest struct {
	Host string `json:"host"`
}
