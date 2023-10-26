package main

import (
	"ece428_mp2/internal/logutil"
	network2 "ece428_mp2/internal/network"
	"ece428_mp2/internal/network/code"
	"flag"
	"github.com/sirupsen/logrus"
	"time"
)

var defaultHostList = []string{
	"fa23-cs425-4801.cs.illinois.edu",
	"fa23-cs425-4802.cs.illinois.edu",
	"fa23-cs425-4803.cs.illinois.edu",
	"fa23-cs425-4804.cs.illinois.edu",
	"fa23-cs425-4805.cs.illinois.edu",
	"fa23-cs425-4806.cs.illinois.edu",
	"fa23-cs425-4807.cs.illinois.edu",
	"fa23-cs425-4808.cs.illinois.edu",
	"fa23-cs425-4809.cs.illinois.edu",
	"fa23-cs425-4810.cs.illinois.edu",
}

func main() {
	serverList := defaultHostList
	var command string
	var target int
	flag.StringVar(&command, "c", "", "determine command name")
	flag.IntVar(&target, "t", 0, "determine command target")
	flag.Parse()

	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	client := network2.NewCallUDPClient()

	if command == "list_mem" {
		r := code.ListMemberRequest{Host: "localhost"}
		req := &network2.CallRequest{
			MethodName: code.ListMember,
			Request:    r,
			TargetHost: serverList[target-1],
		}
		err = client.Call(req)
		if err != nil {
			panic(err)
		}
	} else if command == "list_self" {
		r := code.ListSelfRequest{Host: network2.GetSelfHostName()}
		req := &network2.CallRequest{
			MethodName: code.ListSelf,
			Request:    r,
			TargetHost: serverList[target-1],
		}
		err = client.Call(req)
		if err != nil {
			panic(err)
		}
	} else if command == "leave" {
		r := code.LeaveRequest{Host: network2.GetSelfHostName()}
		req := &network2.CallRequest{
			MethodName: code.Leave,
			Request:    r,
			TargetHost: serverList[target-1],
		}
		err = client.Call(req)
		if err != nil {
			panic(err)
		}
	} else if command == "enable_suspicion" {
		for _, server := range serverList {
			r := code.ChangeSuspicionRequest{SuspicionFlag: true,
				Timestamp: time.Now()}
			req := &network2.CallRequest{
				MethodName: code.ChangeSuspicion,
				Request:    r,
				TargetHost: server,
			}
			err = client.Call(req)
			if err != nil {
				panic(err)
			}
		}
	} else if command == "disable_suspicion" {
		for _, server := range serverList {
			r := code.ChangeSuspicionRequest{SuspicionFlag: false,
				Timestamp: time.Now()}
			req := &network2.CallRequest{
				MethodName: code.ChangeSuspicion,
				Request:    r,
				TargetHost: server,
			}
			err = client.Call(req)
			if err != nil {
				panic(err)
			}
		}
	}
}
