package inpu

import (
	"net/http"
	"slices"
)

type Response struct {
	r *http.Response
}

func newResponse(r *http.Response) *Response {
	return &Response{
		r: r,
	}
}

func (r *Response) IsSuccess() bool {
	return r.r.StatusCode >= 200 && r.r.StatusCode < 300
}

func (r *Response) IsInformational() bool {
	return r.r.StatusCode < 200
}

func (r *Response) IsRedirection() bool {
	return r.r.StatusCode >= 300 && r.r.StatusCode < 400
}

func (r *Response) IsClientError() bool {
	return r.r.StatusCode >= 400 && r.r.StatusCode < 500
}

func (r *Response) IsServerError() bool {
	return r.r.StatusCode >= 500
}

func (r *Response) Is(statusCode int) bool {
	return r.r.StatusCode == statusCode
}

func (r *Response) Status() int {
	return r.r.StatusCode
}

func (r *Response) HttpResponse() *http.Response {
	return r.r
}

func (r *Response) IsOneOf(statusCodes ...int) bool {
	return slices.Contains(statusCodes, r.Status())
}

func (r *Response) ParseJson(t any) error {
	return nil
}

func (r *Response) ParseText() (string, error) {
	return "", nil
}
