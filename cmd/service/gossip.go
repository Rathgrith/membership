package main

import (
	"ece428_mp2"
	"ece428_mp2/internal/logutil"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	service, err := ece428_mp2.NewDefaultGossipGMService()
	service.JoinToGroup([]string{"fa23-cs425-4801.cs.illinois.edu"})
	if err != nil {
		panic(err)
	}

	failListenChan := make(chan string, 10)
	service.SubscribeFailNotification(nil, true, failListenChan)
	go func() {
		for {
			failedHost := <-failListenChan
			fmt.Println(fmt.Sprintf("%v failed", failedHost))
		}
	}()
	service.Serve()
}
