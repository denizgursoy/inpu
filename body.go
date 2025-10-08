package inpu

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	return &xmlBody{
		body: body,
	}
}

func BodyJson(body any) Requester {
	return &jsonBody{
		body: body,
	}
}

func BodyReader(body io.Reader) Requester {
	return &readerBody{body: body}
}

type Requester interface {
	GetBody() (io.Reader, error)
}

type xmlBody struct {
	body any
}

func (x *xmlBody) GetBody() (io.Reader, error) {
	xmlData, err := xml.Marshal(x.body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(xmlData), nil
}

type jsonBody struct {
	body any
}

func (j *jsonBody) GetBody() (io.Reader, error) {
	xmlData, err := json.Marshal(j.body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(xmlData), nil
}

type readerBody struct {
	body io.Reader
}

func (r *readerBody) GetBody() (io.Reader, error) {
	return r.body, nil
}
