package main

import (
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"
	"flag"
	"time"

	"github.com/sirupsen/logrus"
)

var suspicionFlag bool

func init() {
	flag.BoolVar(&suspicionFlag, "suspicion", true, "Set the suspicion flag (true/false)")
}

func main() {
	flag.Parse()

	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	client := network.NewCallUDPClient()
	r := code.ChangeSuspicionRequest{SuspicionFlag: suspicionFlag,
		Timestamp: time.Now()}
	req := &network.CallRequest{
		MethodName: code.ChangeSuspicion,
		Request:    r,
		TargetHost: "localhost",
	}
	err = client.Call(req)
	if err != nil {
		panic(err)
	}
}
