package main

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"
	"flag"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	config.MustLoadGossipFDConfig()
	serverList := config.GetServerList()
	var command string
	var target int
	flag.StringVar(&command, "c", "", "determine command name")
	flag.IntVar(&target, "t", 0, "determine command target")
	flag.Parse()

	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	client := network.NewCallUDPClient()

	if command == "list_mem" {
		r := code.ListMemberRequest{Host: "localhost"}
		req := &network.CallRequest{
			MethodName: code.ListMember,
			Request:    r,
			TargetHost: serverList[target-1],
		}
		err = client.Call(req)
		if err != nil {
			panic(err)
		}
	} else if command == "list_self" {
		r := code.ListSelfRequest{Host: config.GetSelfHostName()}
		req := &network.CallRequest{
			MethodName: code.ListSelf,
			Request:    r,
			TargetHost: serverList[target-1],
		}
		err = client.Call(req)
		if err != nil {
			panic(err)
		}
	} else if command == "leave" {
		r := code.LeaveRequest{Host: config.GetSelfHostName()}
		req := &network.CallRequest{
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
			req := &network.CallRequest{
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
			req := &network.CallRequest{
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
