package main

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg"
	"ece428_mp2/pkg/network"
	"encoding/json"
	"fmt"
	"time"
)

func joinRequest(HostID int, RequestType string, RequestOutTime time.Time) {
	// Create a JoinRequest struct
	request := pkg.JoinRequest{
		HostID:         HostID,
		RequestType:    RequestType,
		RequestOutTime: RequestOutTime,
	}

	// Serialize the struct
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling JoinRequest to JSON:", err)
		return
	}
	// for loop to send 10 requests every 1 second
	for {
		time.Sleep(1 * time.Second)
		// Send serialized data via UDP
		destAddr := "localhost:8000" // Replace with appropriate address and port
		err = network.SendUDP(jsonData, destAddr)
		if err != nil {
			fmt.Println("Error sending UDP request:", err)
			return
		}
		fmt.Println("JoinRequest sent!")
	}
}

func main() {
	host, err := network.GetHostname()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
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
	go joinRequest(id, "join", time.Now())
	// Wait indefinitely so the main function does not exit prematurely
	select {}
}
