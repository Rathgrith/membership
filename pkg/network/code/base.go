package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/*
Format of package:

		| Meta{MagiNumber: xxx, CodeHandlerType: xxx, Length xxx} | RequestHeader | Body  | RequestHeader | Body |
	    |           encode/decode by json                         |     encode/decode according to type in meta  |
*/

const MagicNumber = 0x5d6f

type Meta struct {
	MagicNumber     int32       `json:"magic_number"`
	CodeHandlerType HandlerType `json:"code_handler_type"`
	PayloadLength   uint32      `json:"length"`
}

type Request struct {
	Header *RequestHeader
	Body   RequestBody
}

type RequestHeader struct {
	Method     string `json:"method"` // use this field to indicate the receiver which function will be called
	BodyLength int32  `json:"body_length"`
}

type RequestBody interface{}

/*Handler define the interface to decode/encode message*/
type Handler interface {
	Read(buffer *bytes.Buffer) (*RequestHeader, RequestBody, error)
	Encode(header *RequestHeader, body RequestBody) ([]byte, error)
}

type HandlerInitFunc func() Handler

type HandlerType int32

const (
	JSONType HandlerType = 1
)

var HandlerMap = map[HandlerType]Handler{
	JSONType: NewJSONCodeHandler(),
}

func ReadMetaInfo(buffer *bytes.Buffer) (*Meta, error) {
	meta := Meta{}
	err := binary.Read(buffer, binary.BigEndian, &meta)
	if err != nil {
		return nil, fmt.Errorf("parse meta failed:%w", err)
	}

	return &meta, nil
}
