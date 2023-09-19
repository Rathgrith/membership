package network

import (
	// "gopkg.in/yaml.v2"
	"ece428_mp2/config"
	"ece428_mp2/pkg"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func init() {
	// Initialize membership list
	selfHost, err := GetHostname()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	selfID, err := config.GetHostID(selfHost)
	membershipList[selfID] = pkg.MemberInfo{
		Counter:    1,
		LocalTime:  time.Now(),
		StatusCode: 1,
	}
}

var membershipList = map[int]pkg.MemberInfo{}
var membershipListLock sync.RWMutex

// check if the received udp packet type is join or leave
// if join, add the host to the membership list
// if leave, remove the host from the membership list

func joinMemberToMembershipList(request pkg.JoinRequest, addr net.Addr) {
	membershipListLock.Lock()
	defer membershipListLock.Unlock()
	membershipList[request.HostID] = pkg.MemberInfo{
		Counter:    1,
		LocalTime:  time.Now(),
		StatusCode: 1,
	}
}

func getMembershipList() map[int]pkg.MemberInfo {
	membershipListLock.RLock()
	defer membershipListLock.RUnlock()

	// Return a shallow copy of the membershipList to prevent race conditions
	copiedList := make(map[int]pkg.MemberInfo)
	for k, v := range membershipList {
		copiedList[k] = v
	}
	return copiedList
}

func SendUDPRoutine(HostID int, RequestType string, RequestOutTime time.Time, Destination string) {
	// Create a JoinRequest struct
	request := pkg.JoinRequest{
		HostID:        HostID,
		PacketType:    RequestType,
		PacketOutTime: RequestOutTime,
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
		destAddr := Destination // Replace with appropriate address and port
		fmt.Println("Sending UDP request to", destAddr)
		err = sendUDP(jsonData, destAddr+":8000")
		if err != nil {
			fmt.Println("Error sending UDP request:", err)
			return
		}
		// fmt.Println("JoinRequest sent!")
	}
}

func ReceiveUDPRoutine() {
	selfHost, err := GetHostname()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	selfID, err := config.GetHostID(selfHost)
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
		// fmt.Println("Received", n, "bytes from", addr)
		fmt.Printf("request id: %d, request type: %s, request time: %s\n", request.HostID, request.PacketType, request.PacketOutTime)
		fmt.Println("Received", n, "bytes from", addr)
		if request.PacketType == "joinResponse" {
			// Unmarshal the data and print the data
			var response pkg.JoinResponse
			err = json.Unmarshal(buffer[:n], &response)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return
			}
			for k, v := range response.PacketData {
				fmt.Printf("member id: %d, member counter: %d, member time: %s, member status: %d\n", k, v.Counter, v.LocalTime, v.StatusCode)
			}
		}
		if request.PacketType == "join" {
			joinMemberToMembershipList(request, addr)
			response := pkg.JoinResponse{
				HostID:        selfID,
				PacketType:    "joinResponse",
				PacketOutTime: time.Now(),
				PacketData:    getMembershipList(),
			}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Error marshaling JoinResponse to JSON:", addr.String())
				return
			}
			// Send the response JSON back to the source.
			targetAddr := fmt.Sprintf("%s:8000", addr.(*net.UDPAddr).IP)

			// Send the response JSON back to the target address.
			err = sendUDP(jsonResponse, targetAddr)
			if err != nil {
				fmt.Println("Error sending JoinResponse:", err)
				return
			}
		}
	}
}

func sendUDP(data []byte, destAddr string) error {
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
