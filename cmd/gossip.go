package main

import (
	"ece428_mp2/pkg/gossip"
	"ece428_mp2/pkg/logutil"
	"github.com/sirupsen/logrus"
)

func main() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}
	service := gossip.NewGossipService()
	service.Serve()
}
