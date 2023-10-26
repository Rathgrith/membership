package network

import (
	"bytes"
	code2 "ece428_mp2/internal/network/code"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

const (
	defaultCacheSize = 10
)

type CallRequest struct {
	MethodName code2.MethodType
	Request    interface{}
	TargetHost string
}

type CallUDPClient struct {
	workers      map[string]*udpWorker
	workerPoolMu sync.Mutex
}

var nowaitMethod = map[code2.MethodType]bool{
	code2.ListSelf:        true,
	code2.ListMember:      true,
	code2.Leave:           true,
	code2.ChangeSuspicion: true,
	code2.Heartbeat:       true,
}

func NewCallUDPClient() *CallUDPClient {
	return &CallUDPClient{
		workers: make(map[string]*udpWorker, defaultCacheSize),
	}
}

func (c *CallUDPClient) Call(req *CallRequest) error {
	return c.call(req, req.TargetHost, true)
}

func (c *CallUDPClient) call(req *CallRequest, targetHost string, nowait bool) error {
	encodedReq, err := c.encodeReq(req)
	if err != nil {
		return fmt.Errorf("error arise when encode req:%w", err)
	}

	worker, err := c.getUDPWorker(targetHost)
	if err != nil {
		return fmt.Errorf("error arise when send req:%w", err)
	}

	err = worker.writeRequest(encodedReq)
	if err != nil {
		return fmt.Errorf("error arise when send req:%w", err)
	}

	if !nowait {
		return nil
	}

	return worker.sendRequest()
}

func (c *CallUDPClient) encodeReq(req *CallRequest) ([]byte, error) {
	header := code2.RequestHeader{
		Method: req.MethodName,
	}

	body := req.Request

	encodedData, err := code2.HandlerMap[code2.JSONType].Encode(&header, body)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

func (c *CallUDPClient) getUDPWorker(targetHost string) (*udpWorker, error) {
	_, ok := c.workers[targetHost]
	if !ok {
		c.workerPoolMu.Lock()
		worker, err := newUDPWorker(targetHost)
		if err != nil {
			c.workerPoolMu.Unlock()
			return nil, err
		}
		c.workers[targetHost] = worker
		c.workerPoolMu.Unlock()
	}

	return c.workers[targetHost], nil
}

type udpWorker struct {
	reqBuf *bytes.Buffer // cache encoded request(header + body) for each target (util heartbeat or timeout)
	conn   *net.UDPConn
	mu     sync.Mutex
}

func newUDPWorker(targetHost string) (*udpWorker, error) {
	conn, err := NewUDPConnection(targetHost, DefaultPort)
	if err != nil {
		return nil, fmt.Errorf("build udp conn failed:%w", err)
	}
	worker := &udpWorker{
		conn: conn,
	}

	return worker, nil
}

func (w *udpWorker) initBuf() {
	if w.reqBuf == nil {
		w.reqBuf = bytes.NewBuffer(nil)
	}

	meta := code2.NewMeta(code2.JSONType)
	code2.WriteMeta(meta, w.reqBuf)
}

func (w *udpWorker) writeRequest(encodedReq []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.reqBuf == nil || w.reqBuf.Len() == 0 {
		w.initBuf()
	}
	return binary.Write(w.reqBuf, binary.BigEndian, encodedReq)
}

func (w *udpWorker) sendRequest() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.reqBuf == nil {
		return fmt.Errorf("request buffer is empty")
	}

	_, err := w.conn.Write(w.reqBuf.Bytes())
	if err != nil {
		return fmt.Errorf("send request failed:%w", err)
	}

	w.reqBuf.Reset()

	return nil
}
