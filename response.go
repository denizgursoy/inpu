package inpu

import "net/http"

type Response struct {
}

func newResponse(r http.Response) *Response {
	return &Response{}
}

func (r *Response) IsSuccess() bool {

}

func (r *Response) IsInformational() bool {

}

func (r *Response) IsRedirection() bool {

}

func (r *Response) IsClientError() bool {

}

func (r *Response) IsServerError() bool {

}

func (r *Response) Is(statusCode int) bool {

}

func (r *Response) IsOneOf(statusCode ...int) bool {

}

func (r *Response) ParseJson(t any) error {

}

func (r *Response) ParseText() (string, error) {

}
