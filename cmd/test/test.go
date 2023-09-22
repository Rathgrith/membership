package main

// import (
// 	"ece428_mp2/config"
// 	"ece428_mp2/pkg/gossip"
// 	"fmt"
// 	"time"
// 	// "time"
// )

// var Tfail = time.Second * 2
// var Tcleanup = time.Second * 4

// func main() {
// 	host, err := gossip.GetHostname()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	// GOROUTINE for receiving UDP packets
// 	// GOROUTINE for sending join request
// 	var membershipManager = gossip.NewMembershipManager()
// 	introducer, err := config.GetIntroducer()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	membershipManager.InitMembershiplist(host)
// 	membershipList := membershipManager.GetMembershipList()
// 	gossip.SetMembershipManager(membershipManager)
// 	gossipService := gossip.NewGossipService(membershipManager)
// 	gossipService.Test()
// 	fmt.Println("Membership List:", membershipList)
// 	fmt.Println("Introducer:", introducer)
// 	fmt.Println("Host:", host)
// 	go gossip.ReceiveUDPRoutine()
// 	if introducer != host {
// 		go gossip.SendJoinUDPRoutine(host, "join", introducer)
// 	}
// 	<-gossip.GetJoinCompleteCh() // Wait for the join routine to complete
// 	go membershipManager.StartFailureDetection(Tfail)
// 	go membershipManager.StartCleanupRoutine(Tcleanup)
// 	// Start broadcasting after joining is complete
// 	// gossip.SendSuspicionBroadcast(host, "EnableSuspicionBroadcast")

// 	// Wait indefinitely so the main function does not exit prematurely
// 	select {}
// }
