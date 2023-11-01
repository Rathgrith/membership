package network

import (
	"bytes"
	"context"
	"fmt"
	"membership/internal/network/code"
	"net"
)

const (
	DefaultChanSize       = 10
	DefaultReadBufferSize = 2048
)

type CallUDPServer struct {
	listenPort int
	conn       *net.UDPConn
	f          DispatchFunc

	errChan         chan error
	serveContext    context.Context
	serveCancelFunc context.CancelFunc
}

type DispatchFunc func(header *code.RequestHeader, reqBody []byte) error

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
			s.errChan <- fmt.Errorf("error arise when read from UDP:%s", err.Error())
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
	for dataBuf.Len() > 0 {
		reqHeader, reqBody, err := codeHandler.Read(dataBuf)
		if err != nil {
			continue
		}
		body, ok := reqBody.([]byte)
		if !ok {
			continue
		}
		if int(reqHeader.BodyLength) != len(body) {
			s.errChan <- fmt.Errorf("incomplete request body， method:%v", reqHeader.Method)
			continue
		}
		err = s.f(reqHeader, body)
		if err != nil {
			s.errChan <- fmt.Errorf("handle request failed:%w, string body:%v", err, string(body))
			continue
		}
	}
}

func (s *CallUDPServer) Register(f DispatchFunc) {
	s.f = f
}
