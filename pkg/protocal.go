package pkg

import "time"

type JoinRequest struct {
	HostID         int
	RequestType    string
	RequestOutTime time.Time
}
