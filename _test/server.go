package main

import (
	"ece428_mp2/pkg"
	"ece428_mp2/pkg/gossip"
)

func main() {
	service := gossip.NewGossipService(pkg.NewMembershipManager())
	service.Serve()
}
