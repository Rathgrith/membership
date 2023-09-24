package main

import (
	"ece428_mp2/pkg/logutil"
	"ece428_mp2/pkg/network"
	"ece428_mp2/pkg/network/code"
	"github.com/sirupsen/logrus"
)

type test struct {
	A int    `json:"a"`
	B string `json:"b"`
	C bool   `json:"c"`
}

func main() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}
	client := network.NewCallUDPClient()

	r := code.JoinRequest{Host: "introducer"}
	req := &network.CallRequest{
		MethodName: code.Join,
		Request:    r,
		TargetHost: "127.0.0.1",
	}
	err = client.Call(req)
	if err != nil {
		panic(err)
	}

	err = client.Call(req)
	if err != nil {
		panic(err)
	}

	err = client.Call(req)
	if err != nil {
		panic(err)
	}

	r2 := code.HeartbeatRequest{
		MemberShipList: nil,
	}
	req2 := &network.CallRequest{
		MethodName: code.Heartbeat,
		Request:    r2,
		TargetHost: "127.0.0.1",
	}
	err = client.Call(req2)
	if err != nil {
		panic(err)
	}
}

//func main() {
//	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
//		IP:   net.IPv4(127, 0, 0, 1),
//		Port: 10088,
//	})
//
//	if err != nil {
//		panic(err)
//	}
//
//	wg := sync.WaitGroup{}
//	for i := 0; i < 1; i++ {
//		wg.Add(1)
//		go func() {
//			body := &test{
//				A: 100,
//				B: "Hello World",
//				C: true,
//			}
//			header := &code.RequestHeader{
//				Method: "Test Method",
//			}
//			msg := writeReq(header, body)
//			_, err := conn.Write(msg)
//			if err != nil {
//				panic(err)
//			}
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//}
//
//func writeReq(header *code.RequestHeader, body code.RequestBody) []byte {
//	buf := bytes.Buffer{}
//
//	encodedData, err := code.HandlerMap[code.JSONType].Encode(header, body)
//	if err != nil {
//		panic(err)
//	}
//
//	meta := code.Meta{
//		MagicNumber:     code.MagicNumber,
//		CodeHandlerType: code.JSONType,
//	}
//	err = binary.Write(&buf, binary.BigEndian, meta)
//	if err != nil {
//		panic(err)
//	}
//
//	buf.Write(encodedData)
//	return buf.Bytes()
//}
