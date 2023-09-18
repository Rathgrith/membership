package main

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg/network"
	"fmt"
	"time"
)

func main() {
	host, err := network.GetHostname()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// comment below when on VM
	host = "fa23-cs425-4810.cs.illinois.edu"
	id, err := config.GetHostID(host)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("ID:", id)

	// GOROUTINE for receiving UDP packets
	go network.ReceiveUDPRoutine()
	// GOROUTINE for sending join request
	go network.SendUDPRoutine(id, "join", time.Now())
	// Wait indefinitely so the main function does not exit prematurely
	select {}
}
