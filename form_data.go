package inpu

import (
	"io"
	"net/url"
	"strings"
)

func FormDataFromUrl(v url.Values) io.Reader {
	return strings.NewReader(v.Encode())
}

func FormDataFromMap(v map[string]string) io.Reader {
	values := url.Values{}
	for key, val := range v {
		values.Set(key, val)
	}

	return strings.NewReader(values.Encode())
}

func FormDataFrom(v map[string][]string) io.Reader {
	return strings.NewReader(url.Values(v).Encode())
}
