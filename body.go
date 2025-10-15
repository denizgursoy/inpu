package inpu

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"strings"
)

func BodyFormDataFromMap(body map[string]string) Requester {
	values := url.Values{}
	for key, val := range body {
		values.Set(key, val)
	}

	return BodyFormData(values)
}

func BodyFormData(body map[string][]string) Requester {
	return BodyReader(strings.NewReader(url.Values(body).Encode()))
}

func BodyString(body string) Requester {
	return BodyReader(strings.NewReader(body))
}

func BodyXml(body any) Requester {
	xmlData, err := xml.Marshal(body)
	if err != nil {
		return newRequestBody(nil, fmt.Errorf("could not marshal to XML: %w", err))
	}

	return newRequestBody(bytes.NewBuffer(xmlData), nil)
}

func BodyJson(body any) Requester {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return newRequestBody(nil, fmt.Errorf("could not marshal to JSON: %w", err))
	}

	return newRequestBody(bytes.NewBuffer(jsonData), nil)
}

func BodyReader(body io.Reader) Requester {
	return newRequestBody(body, nil)
}

type Requester interface {
	GetBody() (io.Reader, error)
}

type requestBody struct {
	body io.Reader
	err  error
}

func newRequestBody(body io.Reader, err error) *requestBody {
	return &requestBody{body: body, err: err}
}

func (r *requestBody) GetBody() (io.Reader, error) {
	return r.body, r.err
}
