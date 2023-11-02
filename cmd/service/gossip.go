package main

import (
	"fmt"
	"github.com/Rathgrith/membership"
)

func main() {
	service, err := membership.NewDefaultGossipGMService()
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
