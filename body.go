package inpu

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/url"
	"strings"
)

func BodyFormDataFromMap(body map[string]string) io.Reader {
	values := url.Values{}
	for key, val := range body {
		values.Set(key, val)
	}

	return BodyFormData(values)
}

func BodyFormData(body map[string][]string) io.Reader {
	return strings.NewReader(url.Values(body).Encode())
}

func BodyString(body string) io.Reader {
	return strings.NewReader(body)
}

func BodyXml(body any) Requester {
	return &requestBody{
		body: body,
	}
}

type Requester interface {
	GetBody() (io.Reader, error)
}

type requestBody struct {
	body any
}

func (r *requestBody) GetBody() (io.Reader, error) {
	xmlData, err := xml.Marshal(r.body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(xmlData), nil
}
