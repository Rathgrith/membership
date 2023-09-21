package main

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg"
	"ece428_mp2/pkg/network"
	"fmt"
	// "time"
)

func main() {
	host, err := network.GetHostname()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// GOROUTINE for receiving UDP packets
	// GOROUTINE for sending join request
	introducer, err := config.GetIntroducer()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	pkg.InitMembershiplist(host)
	membershipList := pkg.GetMembershipList()
	fmt.Println("Membership List:", membershipList)
	fmt.Println("Introducer:", introducer)
	fmt.Println("Host:", host)
	go network.ReceiveUDPRoutine()
	if introducer != host {
		go network.SendJoinUDPRoutine(host, "join", introducer)
	}
	<-network.GetJoinCompleteCh() // Wait for the join routine to complete

	// Start broadcasting after joining is complete
	network.SendSuspicionBroadcast(host, "EnableSuspicionBroadcast")

	// Wait indefinitely so the main function does not exit prematurely
	select {}
}
