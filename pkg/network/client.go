package network

import "ece428_mp2/pkg/network/code"

type CallRequest struct {
	MethodName   string
	Request      interface{}
	TargetMember string
}

type CallUDPClient struct {
}

func (c *CallUDPClient) Call(req *CallRequest) {
	header := code.RequestHeader{
		Method: req.MethodName,
	}
	body := req.Request

}
