package inpu

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
	"strconv"
	"time"
)

type Req struct {
	userClient           *http.Client
	genericTimeout       time.Duration
	httpReq              *http.Request
	requestCreationError error
	bodyCreationError    error
}

func Get(url string) *Req {
	return getReq(context.Background(), url, nil, nil, nil, 0)
}

func GetCtx(ctx context.Context, url string) *Req {
	return getReq(ctx, url, nil, nil, nil, 0)
}

func Post(url string, body any) *Req {
	return postReq(context.Background(), url, body, nil, nil, nil, 0)
}

func PostCtx(ctx context.Context, url string, body any) *Req {
	return postReq(ctx, url, body, nil, nil, nil, 0)
}

func Delete(url string, body any) *Req {
	return deleteReq(context.Background(), url, body, nil, nil, nil, 0)
}

func DeleteCtx(ctx context.Context, url string, body any) *Req {
	return deleteReq(ctx, url, body, nil, nil, nil, 0)
}

func Put(url string, body any) *Req {
	return putReq(context.Background(), url, body, nil, nil, nil, 0)
}

func PutCtx(ctx context.Context, url string, body any) *Req {
	return putReq(ctx, url, body, nil, nil, nil, 0)
}

func Patch(url string, body any) *Req {
	return patchReq(context.Background(), url, body, nil, nil, nil, 0)
}

func PatchCtx(ctx context.Context, url string, body any) *Req {
	return patchReq(ctx, url, body, nil, nil, nil, 0)
}

func getReq(ctx context.Context, url string, headers http.Header, queries netUrl.Values, client *http.Client, genericTimeout time.Duration) *Req {
	return newRequest(ctx, http.MethodGet, url, nil, headers, queries, client, genericTimeout)
}

func postReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values, client *http.Client, genericTimeout time.Duration) *Req {
	return newRequest(ctx, http.MethodPost, url, body, headers, queries, client, genericTimeout)
}

func deleteReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values, client *http.Client, genericTimeout time.Duration) *Req {
	return newRequest(ctx, http.MethodDelete, url, body, headers, queries, client, genericTimeout)
}

func putReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values, client *http.Client, genericTimeout time.Duration) *Req {
	return newRequest(ctx, http.MethodPut, url, body, headers, queries, client, genericTimeout)
}

func patchReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values, client *http.Client, genericTimeout time.Duration) *Req {
	return newRequest(ctx, http.MethodPatch, url, body, headers, queries, client, genericTimeout)
}

func newRequest(ctx context.Context, method, rawUrl string, body any, headers http.Header, queries netUrl.Values, userClient *http.Client, genericTimeout time.Duration) *Req {
	reader, bodyCreationError := getBody(body)
	httpReq, requestCreationError := http.NewRequestWithContext(ctx, method, rawUrl, reader)
	isSuccessful := requestCreationError == nil && bodyCreationError == nil
	if headers != nil && isSuccessful {
		for k, v := range headers {
			for _, v1 := range v {
				httpReq.Header.Add(k, v1)
			}
		}
	}
	if queries != nil && isSuccessful {
		query := httpReq.URL.Query()
		for k, v := range queries {
			for _, v1 := range v {
				query.Add(k, v1)
			}
		}
		httpReq.URL.RawQuery = query.Encode()
	}

	return &Req{
		userClient:           userClient,
		genericTimeout:       genericTimeout,
		requestCreationError: requestCreationError,
		bodyCreationError:    bodyCreationError,
		httpReq:              httpReq,
	}
}

func (r *Req) isSuccessfullyCreated() bool {
	return r.requestCreationError == nil && r.bodyCreationError == nil
}

func (r *Req) Header(key, val string) *Req {
	r.addHeader(key, val)
	return r
}

func (r *Req) ContentTypeJson() *Req {
	r.ContentType(MimeTypeJson)

	return r
}

func (r *Req) ContentTypeText() *Req {
	r.ContentType(MimeTypeText)

	return r
}

func (r *Req) ContentTypeHtml() *Req {
	r.ContentType(MimeTypeHtml)

	return r
}

func (r *Req) ContentType(contentType string) *Req {
	r.addHeader(HeaderContentType, contentType)

	return r
}

func (r *Req) AuthBasic(username, password string) *Req {
	cred := username + ":" + password
	r.addHeader(HeaderAuthorization, "Basic "+base64.StdEncoding.EncodeToString([]byte(cred)))

	return r
}

func (r *Req) AuthToken(token string) *Req {
	r.addHeader(HeaderAuthorization, "Bearer "+token)
	return r
}

func (r *Req) AcceptJson() *Req {
	r.addHeader(HeaderAccept, MimeTypeJson)
	return r
}

func (r *Req) addQueryValue(key, value string) *Req {
	if r.isSuccessfullyCreated() {
		query := r.httpReq.URL.Query()
		query.Add(key, value)
		r.httpReq.URL.RawQuery = query.Encode()
	}
	return r
}

func (r *Req) addHeader(key, value string) *Req {
	if r.isSuccessfullyCreated() {
		r.httpReq.Header.Add(key, value)
	}

	return r
}

func (r *Req) QueryInt8(name string, v int8) *Req {
	return r.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (r *Req) QueryInt16(name string, v int16) *Req {
	return r.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (r *Req) QueryInt32(name string, v int32) *Req {
	return r.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (r *Req) QueryInt(name string, v int) *Req {
	return r.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (r *Req) QueryInt64(name string, v int64) *Req {
	return r.addQueryValue(name, strconv.FormatInt(v, 10))
}

func (r *Req) QueryUint8(name string, v uint8) *Req {
	return r.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (r *Req) QueryUint16(name string, v uint16) *Req {
	return r.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (r *Req) QueryUint32(name string, v uint32) *Req {
	return r.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (r *Req) QueryUint(name string, v uint) *Req {
	return r.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (r *Req) QueryUint64(name string, v uint64) *Req {
	return r.addQueryValue(name, strconv.FormatUint(v, 10))
}

func (r *Req) QueryFloat32(name string, v float32) *Req {
	return r.addQueryValue(name, strconv.FormatFloat(float64(v), 'f', -1, 64))
}

func (r *Req) QueryFloat64(name string, v float64) *Req {
	return r.addQueryValue(name, strconv.FormatFloat(v, 'f', -1, 64))
}

func (r *Req) QueryBool(name string, v bool) *Req {
	return r.addQueryValue(name, strconv.FormatBool(v))
}

func (r *Req) QueryString(name string, v string) *Req {
	return r.addQueryValue(name, v)
}

func (r *Req) QueryInt8Ptr(name string, v *int8) *Req {
	if v == nil {
		return r
	}
	return r.QueryInt8(name, *v)
}

func (r *Req) QueryInt16Ptr(name string, v *int16) *Req {
	if v == nil {
		return r
	}
	return r.QueryInt16(name, *v)
}

func (r *Req) QueryInt32Ptr(name string, v *int32) *Req {
	if v == nil {
		return r
	}
	return r.QueryInt32(name, *v)
}

func (r *Req) QueryIntPtr(name string, v *int) *Req {
	if v == nil {
		return r
	}
	return r.QueryInt(name, *v)
}

func (r *Req) QueryInt64Ptr(name string, v *int64) *Req {
	if v == nil {
		return r
	}
	return r.QueryInt64(name, *v)
}

func (r *Req) QueryUint8Ptr(name string, v *uint8) *Req {
	if v == nil {
		return r
	}
	return r.QueryUint8(name, *v)
}

func (r *Req) QueryUint16Ptr(name string, v *uint16) *Req {
	if v == nil {
		return r
	}
	return r.QueryUint16(name, *v)
}

func (r *Req) QueryUint32Ptr(name string, v *uint32) *Req {
	if v == nil {
		return r
	}
	return r.QueryUint32(name, *v)
}

func (r *Req) QueryUintPtr(name string, v *uint) *Req {
	if v == nil {
		return r
	}
	return r.QueryUint(name, *v)
}

func (r *Req) QueryUint64Ptr(name string, v *uint64) *Req {
	if v == nil {
		return r
	}
	return r.QueryUint64(name, *v)
}

func (r *Req) QueryFloat32Ptr(name string, v *float32) *Req {
	if v == nil {
		return r
	}
	return r.QueryFloat32(name, *v)
}

func (r *Req) QueryFloat64Ptr(name string, v *float64) *Req {
	if v == nil {
		return r
	}
	return r.QueryFloat64(name, *v)
}

func (r *Req) QueryBoolPtr(name string, v *bool) *Req {
	if v == nil {
		return r
	}
	return r.QueryBool(name, *v)
}

func (r *Req) QueryStringPtr(name string, v *string) *Req {
	if v == nil {
		return r
	}
	return r.QueryString(name, *v)
}

func (r *Req) Send() (*Response, error) {
	if !r.isSuccessfullyCreated() {
		return nil, errors.Join(r.requestCreationError, r.bodyCreationError)
	}
	client := http.DefaultClient
	if r.userClient != nil {
		client = r.userClient
	}

	httpResponse, err := client.Do(r.httpReq)
	if err != nil {
		return nil, fmt.Errorf("could not send the request: %w,%w", ErrConnectionFailed, err)
	}

	return newResponse(httpResponse), nil
}

func getBody(reqBody any) (io.Reader, error) {
	var body io.Reader
	if reqBody != nil {
		switch v := reqBody.(type) {
		case io.Reader:
			body = v
		default:
			bodyAsBytes, err := json.Marshal(reqBody)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(bodyAsBytes)
		}
	}

	return body, nil
}
