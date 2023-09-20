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
	// comment below when on VM
	// host = "fa23-cs425-4810.cs.illinois.edu"
	// id, err := config.GetHostID(host)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("ID:", id)

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
	if introducer != host {
		go network.SendJoinUDPRoutine(host, "join", introducer)
	}
	go network.ReceiveUDPRoutine()
	// Wait indefinitely so the main function does not exit prematurely
	select {}
}
