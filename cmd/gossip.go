package main

import (
	"ece428_mp2/config"
	"ece428_mp2/pkg/gossipGM"
	"ece428_mp2/pkg/logutil"

	"github.com/sirupsen/logrus"
)

func main() {
	config.MustLoadGossipFDConfig()
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	service := gossipGM.NewGossipService()
	service.Serve()
}
