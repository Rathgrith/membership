package network

import (
	// "gopkg.in/yaml.v2"

	"ece428_mp2/pkg"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var membershipList = map[int]pkg.MemberInfo{}
var stopSendJoinCh = make(chan struct{})
var joinCompleteCh = make(chan struct{})
var closeOnce sync.Once
var Port = ":8000"
var BufferLen = 1024

func SendJoinUDPRoutine(Host string, RequestType string, Destination string) {
	// Create a JoinRequest struct
	request := pkg.JoinRequest{
		Host:       Host,
		PacketType: RequestType,
	}

	// Serialize the struct
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling JoinRequest to JSON:", err)
		return
	}

	for {
		select {
		case <-stopSendJoinCh:
			return
		default:

			// Send serialized data via UDP
			destAddr := Destination // Replace with appropriate address and port
			fmt.Println("Sending UDP request to", destAddr)
			err = SendUDP(jsonData, destAddr+":8000")
			if err != nil {
				fmt.Println("Error sending UDP request:", err)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// ReceiveUDPRoutine listens for incoming UDP packets and processes them
func ReceiveUDPRoutine() {
	hostname, err := GetHostname()
	if err != nil {
		log.Printf("Error getting hostname: %s", err)
		return
	}

	pc, err := net.ListenPacket("udp", Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", Port, err)
	}
	defer pc.Close()

	buffer := make([]byte, BufferLen)
	for {
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading UDP packet: %s", err)
			return
		}

		go handleIncomingPacket(buffer[:n], addr, hostname)
	}
}

func SendSuspicionBroadcast(Host string, BroadcastType string) {
	if BroadcastType != "EnableSuspicionBroadcast" && BroadcastType != "DisableSuspicionBroadcast" {
		log.Printf("Invalid BroadcastType: %s", BroadcastType)
		return
	}

	broadcast := pkg.Broadcast{
		Host:         Host,
		PacketType:   BroadcastType,
		BroadcastTTL: 3, // Set the TTL value as per your requirements
	}

	// Serialize the Broadcast struct
	jsonData, err := json.Marshal(broadcast)
	if err != nil {
		log.Printf("Error marshaling Broadcast to JSON: %s", err)
		return
	}

	// Send to all nodes in the membership list
	for _, memberInfo := range pkg.GetMembershipList() {
		// Don't send back to ourselves
		if memberInfo.Hostname != Host {
			targetAddr := memberInfo.Hostname + Port
			if err := SendUDP(jsonData, targetAddr); err != nil {
				log.Printf("Error sending Broadcast to %s: %s", targetAddr, err)
			}
		}
	}
}

func handleIncomingPacket(data []byte, addr net.Addr, selfHost string) {
	var request pkg.JoinRequest
	if err := json.Unmarshal(data, &request); err != nil {
		log.Printf("Error unmarshaling JSON: %s", err)
		return
	}

	log.Printf("Received request from %s: %s of type %s", addr, request.Host, request.PacketType)
	switch request.PacketType {
	case "join":
		if request.Host != selfHost {
			handleJoinRequest(request, addr, selfHost)
		}
	case "joinResponse":
		handleJoinResponse(data)
	case "EnableSuspicionBroadcast":
		fmt.Printf("Received enable suspicion broadcast from %s", request.Host)
		handleBroadcast(data, selfHost)
	case "DisableSuspicionBroadcast":
		fmt.Printf("Received enable suspicion broadcast from %s", request.Host)
		handleBroadcast(data, selfHost)
	default:
		log.Printf("Unknown packet type %s", request.PacketType)
	}
}

func handleJoinRequest(request pkg.JoinRequest, addr net.Addr, selfHost string) {
	pkg.JoinToMembershipList(request, request.Host)

	response := pkg.JoinResponse{
		Host:          selfHost,
		PacketType:    "joinResponse",
		PacketOutTime: time.Now(),
		PacketData:    pkg.GetMembershipList(),
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling JoinResponse: %s", err)
		return
	}

	targetAddr := addr.(*net.UDPAddr).IP.String() + Port
	if err := SendUDP(data, targetAddr); err != nil {
		log.Printf("Error sending JoinResponse: %s", err)
	}
}

func handleBroadcast(data []byte, selfHost string) {
	var broadcast pkg.Broadcast
	if err := json.Unmarshal(data, &broadcast); err != nil {
		log.Printf("Error unmarshaling Broadcast: %s", err)
		return
	}

	// Process the broadcast message here (if necessary)
	log.Printf("Received broadcast from %s with type %s and TTL %d",
		broadcast.Host, broadcast.PacketType, broadcast.BroadcastTTL)

	// Decrement TTL
	broadcast.BroadcastTTL--

	// If TTL is still positive, forward the broadcast to other nodes
	if broadcast.BroadcastTTL > 0 {
		fmt.Printf("Forwarding broadcast from %s with type %s and TTL %d", broadcast.Host, broadcast.PacketType, broadcast.BroadcastTTL)
		ForwardBroadcast(broadcast, selfHost)
	}
}

func ForwardBroadcast(broadcast pkg.Broadcast, selfHost string) {
	data, err := json.Marshal(broadcast)
	if err != nil {
		log.Printf("Error marshaling Broadcast: %s", err)
		return
	}

	// Forward the broadcast message to other nodes in the membership list
	for _, memberInfo := range pkg.GetMembershipList() {
		// Don't send back to the source or to ourselves
		if memberInfo.Hostname != selfHost {
			targetAddr := memberInfo.Hostname + Port
			if err := SendUDP(data, targetAddr); err != nil {
				log.Printf("Error sending Broadcast to %s: %s", targetAddr, err)
			}
		}
		// if there is no other node in the membership list, report that and do nothing
		if len(pkg.GetMembershipList()) == 1 {
			fmt.Println("No other node in the membership list")
		}
	}
}

func handleJoinResponse(data []byte) {
	var response pkg.JoinResponse
	if err := json.Unmarshal(data, &response); err != nil {
		log.Printf("Error unmarshaling JoinResponse: %s", err)
		return
	}

	closeOnce.Do(func() {
		close(stopSendJoinCh)
		close(joinCompleteCh)
	})

	pkg.OverwriteMembershipList(response.PacketData)
	log.Println("Membership list updated!")
	for k, v := range pkg.GetMembershipList() {
		log.Printf("Member: %s, Counter: %d, Time: %s, Status: %d, Hostname: %s",
			k, v.Counter, v.LocalTime, v.StatusCode, v.Hostname)
	}
}

func SendUDP(data []byte, destAddr string) error {
	addr, err := net.ResolveUDPAddr("udp", destAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	return err
}

// GetHostname returns the hostname of the machine.
func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
	// hostname format fa23-cs425-48XX.cs.illinois.edu
}

func GetJoinCompleteCh() chan struct{} {
	return joinCompleteCh
}
