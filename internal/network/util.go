package network

import (
	"ece428_mp2/config"
	"fmt"
	"net"
	"os"
)

func CleanUDPReceiveBuffer() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: config.GetListenPort(),
	})
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	var data [1024]byte
	fmt.Println(listen.ReadFromUDP(data[:]))
}

func GetSelfHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
	// hostname format fa23-cs425-48XX.cs.illinois.edu
}
