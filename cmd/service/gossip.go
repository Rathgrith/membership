package main

import (
	"ece428_mp2"
	"ece428_mp2/internal/logutil"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	// log is just for test, should not be used in prod env
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	service, err := ece428_mp2.NewDefaultGossipGMService()
	service.JoinToGroup([]string{"fa23-cs425-4801.cs.illinois.edu"}) // join a group explicitly
	if err != nil {
		panic(err)
	}

	// subscribe member failed event
	failListenChan := make(chan string, 10)
	service.SubscribeFailNotification(nil, true, failListenChan)
	go func() { // handle failed event
		for {
			failedHost := <-failListenChan
			fmt.Println(fmt.Sprintf("%v failed", failedHost))
		}
	}()

	// get all alive hosts
	fmt.Println(service.GetHostsOfAllMembers())

	// start
	service.Serve()
}
