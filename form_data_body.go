package inpu

import (
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
