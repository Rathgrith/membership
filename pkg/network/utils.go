package network

import (
	// "gopkg.in/yaml.v2"
	"ece428_mp2/pkg"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

func ReceiveUDPRoutine() {
	// listen to port 8000 for upcomming UDP packets
	pc, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	// loop to receive packets
	for {
		buffer := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			return
		}
		// unmarshal the data and print the data
		var request pkg.JoinRequest
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
		fmt.Printf("request id: %d, request type: %s, request time: %s\n", request.HostID, request.RequestType, request.RequestOutTime)
		fmt.Println("Received", n, "bytes from", addr)
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

// write me a function to match the hostname to id: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
// Hostname format: fa23-cs425-48XX.cs.illinois.edu
// func GetID() (int, error) {
// 	hostname, err := GetHostname()
// 	if err != nil {
// 		return -1, err
// 	}
// 	fmt.Println("Hostname:", hostname)
// 	fmt.Println("Hostname[11:13]:", hostname[11:13])
// 	//convert string to int
// 	return strconv.Atoi(hostname[11:13])
// }
