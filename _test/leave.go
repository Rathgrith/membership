package main

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"

	"github.com/sirupsen/logrus"
)

func main() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}
	client := network.NewCallUDPClient()
	r := code.LeaveRequest{Host: config.GetSelfHostName()}
	req := &network.CallRequest{
		MethodName: code.Leave,
		Request:    r,
		TargetHost: "localhost",
	}
	err = client.Call(req)
	if err != nil {
		panic(err)
	}
}
