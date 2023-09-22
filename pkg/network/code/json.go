package code

import (
	"bytes"
	"ece428_mp2/pkg/logutil"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type test struct {
	A int    `json:"a"`
	B string `json:"b"`
	C bool   `json:"c"`
}

const (
	headerOffset = 29
)

type JSONCodeHandler struct {
}

func NewJSONCodeHandler() Handler {
	return &JSONCodeHandler{}
}

func (J JSONCodeHandler) Read(buffer *bytes.Buffer) (*RequestHeader, RequestBody, error) {
	header := RequestHeader{}
	rawHeader := buffer.Next(headerOffset)
	err := json.Unmarshal(rawHeader, &header)
	if err != nil {
		logutil.Logger.Error(err)
		return nil, nil, fmt.Errorf("parse req header failed:%w", err)
	}
	logutil.Logger.Infof("header:%v", header)

	body := buffer.Next(int(header.BodyLength))

	_ = Request{
		Header: &header,
		Body:   body,
	}
	return &header, body, nil
}

func (J JSONCodeHandler) Encode(header *RequestHeader, body RequestBody) ([]byte, error) {
	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error raised in the body encode:%w", err)
	}

	header.BodyLength = int32(len(rawBody))
	rawHeader, err := json.Marshal(header)
	fmt.Println(fmt.Sprintf("header len:%v", len(rawHeader)))
	if err != nil {
		return nil, fmt.Errorf("error raised in the header encode:%w", err)
	}

	buf := bytes.Buffer{}
	err = binary.Write(&buf, binary.BigEndian, rawHeader)
	err = binary.Write(&buf, binary.BigEndian, rawBody)

	logutil.Logger.Infof(buf.String())

	return buf.Bytes(), err
}
