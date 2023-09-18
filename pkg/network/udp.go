package network

import (
	// "fmt"
	"net"
)

const (
	connType = "udp"
)

func NewUDPConnection(connHost string, connPort string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr(connType, connHost+":"+connPort)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(connType, nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// func main() {
// 	// Example usage
// 	conn, err := NewUDPConnection("127.0.0.1", "8080")
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	fmt.Println("UDP connection established!")
// }
