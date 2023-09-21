package main

import (
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"github.com/sirupsen/logrus"
)

func main() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}
	server, err := network.NewUDPServer(10088)
	if err != nil {
		logutil.Logger.Error(err)
		panic(err)
	}

	errChan := server.Serve()

	logutil.Logger.Debug("server started!")
	select {
	case err = <-errChan:
		panic(err)
	}
}
