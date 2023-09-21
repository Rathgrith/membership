package network

import (
	"bytes"
	"context"
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network/code"
	"encoding/json"
	"fmt"
	"net"
)

const (
	DefaultChanSize       = 10
	DefaultReadBufferSize = 2048
)

type CallUDPServer struct {
	listenPort int
	conn       *net.UDPConn

	errChan         chan error
	serveContext    context.Context
	serveCancelFunc context.CancelFunc
}

func NewUDPServer(listenPort int) (*CallUDPServer, error) {
	return &CallUDPServer{listenPort: listenPort}, nil
}

func (s *CallUDPServer) Serve() <-chan error {
	s.serveContext, s.serveCancelFunc = context.WithCancel(context.Background())
	errChan := make(chan error, DefaultChanSize)
	s.errChan = errChan
	go s.serve(s.serveContext)

	return errChan
}

func (s *CallUDPServer) serve(ctx context.Context) {
	conn, err := NewUDPListenConn(DefaultIP, s.listenPort)
	if err != nil {
		panic(fmt.Errorf("init UDP listner failed:%w", err))
	}

	data := make([]byte, DefaultReadBufferSize)
	for {
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			logutil.Logger.Errorf("error arise when read from UDP:%s", err.Error())
			s.errChan <- err
			continue
		}

		if n != 0 {
			buf := bytes.NewBuffer(data[:n])
			go s.serveUDPRequest(ctx, buf)
		}
	}
}

func (s *CallUDPServer) serveUDPRequest(ctx context.Context, dataBuf *bytes.Buffer) {
	meta, err := code.ReadMetaInfo(dataBuf)
	if err != nil {
		s.errChan <- fmt.Errorf("read meta info failed:%w", err)
		return
	}

	codeHandler := code.HandlerMap[meta.CodeHandlerType]
	reqHeader, reqBody, err := codeHandler.Read(dataBuf)

	header, err := json.Marshal(reqHeader)
	body, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}
	logutil.Logger.Infof("recieved request; header:%v body:%v", string(header), string(body))
}

func Register(receiver interface{}) error {
	_ = newService(receiver)
	return nil
}
