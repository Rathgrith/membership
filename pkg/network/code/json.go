package code

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type JSONCodeHandler struct {
}

func NewJSONCodeHandler() Handler {
	return &JSONCodeHandler{}
}

func (J JSONCodeHandler) ReadMetaInfo(buffer *bytes.Buffer) (*Meta, error) {
	meta := Meta{}
	err := binary.Read(buffer, binary.BigEndian, &meta)
	if err != nil {
		return nil, fmt.Errorf("parse meta failed:%w", err)
	}

	return &meta, nil
}

func (J JSONCodeHandler) Read(buffer *bytes.Buffer) (*RequestHeader, RequestBody, error) {
	header := RequestHeader{}
	dec := json.NewDecoder(buffer)
	err := dec.Decode(&header)
	if err != nil {
		return nil, nil, fmt.Errorf("parse req header failed:%w", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("error raised in the header encode:%w", err)
	}

	buf := bytes.Buffer{}
	err = binary.Write(&buf, binary.BigEndian, rawHeader)
	err = binary.Write(&buf, binary.BigEndian, rawBody)

	return buf.Bytes(), err
}
