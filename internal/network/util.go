package network

import (
	"os"
)

func GetSelfHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
	// hostname format fa23-cs425-48XX.cs.illinois.edu
}
