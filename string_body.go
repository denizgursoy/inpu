package inpu

import (
	"io"
	"strings"
)

func BodyString(body string) io.Reader {
	return strings.NewReader(body)
}
