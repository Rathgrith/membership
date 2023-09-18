package main

import (
	"ece428_mp2/pkg/network"
	"fmt"
)

func main() {
	host, err := network.GetHostname()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Hostname:", host)

}
