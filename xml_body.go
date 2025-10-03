package inpu

import (
	"bytes"
	"encoding/xml"
	"io"
)

type Requester interface {
	GetBody() (io.Reader, error)
}

type requestBody struct {
	body any
}

func BodyXml(body any) Requester {
	return &requestBody{
		body: body,
	}
}

func (r *requestBody) GetBody() (io.Reader, error) {
	xmlData, err := xml.Marshal(r.body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(xmlData), nil
}
