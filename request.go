package inpu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"
)

type Req struct {
	userClient           *http.Client
	httpReq              *http.Request
	requestCreationError error
	timeOut              time.Duration
}

func Get(url string) *Req {
	return getReq(context.Background(), url, nil, nil, nil, "")
}

func GetCtx(ctx context.Context, url string) *Req {
	return getReq(ctx, url, nil, nil, nil, "")
}

func Post(url string, body any) *Req {
	return postReq(context.Background(), url, body, nil, nil, nil, "")
}

func PostCtx(ctx context.Context, url string, body any) *Req {
	return postReq(ctx, url, body, nil, nil, nil, "")
}

func Delete(url string, body any) *Req {
	return deleteReq(context.Background(), url, body, nil, nil, nil, "")
}

func DeleteCtx(ctx context.Context, url string, body any) *Req {
	return deleteReq(ctx, url, body, nil, nil, nil, "")
}

func Put(url string, body any) *Req {
	return putReq(context.Background(), url, body, nil, nil, nil, "")
}

func PutCtx(ctx context.Context, url string, body any) *Req {
	return putReq(ctx, url, body, nil, nil, nil, "")
}

func Patch(url string, body any) *Req {
	return patchReq(context.Background(), url, body, nil, nil, nil, "")
}

func PatchCtx(ctx context.Context, url string, body any) *Req {
	return patchReq(ctx, url, body, nil, nil, nil, "")
}

func getReq(ctx context.Context, url string, headers http.Header, queries netUrl.Values,
	client *http.Client, path string,
) *Req {
	return newRequest(ctx, http.MethodGet, url, nil, headers, queries, client, path)
}

func postReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values,
	client *http.Client, path string,
) *Req {
	return newRequest(ctx, http.MethodPost, url, body, headers, queries, client, path)
}

func deleteReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values,
	client *http.Client, path string,
) *Req {
	return newRequest(ctx, http.MethodDelete, url, body, headers, queries, client, path)
}

func putReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values,
	client *http.Client, path string,
) *Req {
	return newRequest(ctx, http.MethodPut, url, body, headers, queries, client, path)
}

func patchReq(ctx context.Context, url string, body any, headers http.Header, queries netUrl.Values,
	client *http.Client, path string,
) *Req {
	return newRequest(ctx, http.MethodPatch, url, body, headers, queries, client, path)
}

func newRequest(ctx context.Context, method, path string, body any, headers http.Header, queries netUrl.Values,
	userClient *http.Client, basePath string,
) *Req {
	bodyAsReader, err := getBody(body)
	if err != nil {
		return newInvalidRequest(fmt.Errorf("%w: %w", ErrInvalidBody, err))
	}
	url, err := getUrl(basePath, path)
	if err != nil {
		return newInvalidRequest(err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url.String(), bodyAsReader)
	if err != nil {
		return newInvalidRequest(fmt.Errorf("%w: %w", ErrRequestCreationFailed, err))
	}

	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			for _, v1 := range v {
				httpReq.Header.Add(k, v1)
			}
		}
	}
	if queries != nil && len(queries) > 0 {
		query := httpReq.URL.Query()
		for k, v := range queries {
			for _, v1 := range v {
				query.Add(k, v1)
			}
		}
		httpReq.URL.RawQuery = query.Encode()
	}

	return &Req{
		userClient: userClient,
		httpReq:    httpReq,
	}
}

func newInvalidRequest(err error) *Req {
	return &Req{
		requestCreationError: err,
	}
}

func getUrl(basePath, path string) (*netUrl.URL, error) {
	ref, err := netUrl.Parse(path)
	if err != nil {
		return ref, fmt.Errorf("%w: %w", ErrCouldNotParsePath, err)
	}

	if len(strings.TrimSpace(basePath)) > 0 {
		base, err := netUrl.Parse(basePath)
		if err != nil {
			return base, fmt.Errorf("%w: %w", ErrCouldNotParseBaseUrl, err)
		}

		return base.ResolveReference(ref), nil
	}

	return ref, nil
}

func (r *Req) isSuccessfullyCreated() bool {
	return r.requestCreationError == nil
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

func (r *Req) ContentTypeXml() *Req {
	r.ContentType(MimeTypeApplicationXml)

	return r
}

func (r *Req) ContentTypeFormUrlEncoded() *Req {
	r.ContentType(MimeTypeFormUrlEncoded)

	return r
}

func (r *Req) ContentType(contentType string) *Req {
	r.addHeader(HeaderContentType, contentType)

	return r
}

func (r *Req) AuthBasic(username, password string) *Req {
	r.addHeader(HeaderAuthorization, getBasicAuthHeaderValue(username, password))

	return r
}

func (r *Req) AuthToken(token string) *Req {
	r.addHeader(HeaderAuthorization, getTokenHeaderValue(token))
	return r
}

func (r *Req) AcceptJson() *Req {
	r.addHeader(HeaderAccept, MimeTypeJson)

	return r
}

func (r *Req) UserAgent(userAgent string) *Req {
	r.addHeader(HeaderUserAgent, userAgent)

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

func (r *Req) TimeOutIn(duration time.Duration) *Req {
	r.timeOut = duration

	return r
}

func (r *Req) Send() (*Response, error) {
	if !r.isSuccessfullyCreated() {
		return nil, r.requestCreationError
	}
	client := getDefaultClient()
	if r.userClient != nil {
		client = r.userClient
	}
	if r.timeOut > 0 {
		timeoutCtx, cancel := context.WithTimeout(r.httpReq.Context(), r.timeOut)
		defer cancel()

		r.httpReq = r.httpReq.WithContext(timeoutCtx)
	}

	httpResponse, err := client.Do(r.httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w,%w", ErrConnectionFailed, err)
	}

	return newResponse(httpResponse), nil
}

func getBody(reqBody any) (io.Reader, error) {
	var body io.Reader
	if reqBody != nil {
		switch v := reqBody.(type) {
		case io.Reader:
			body = v
		case Requester:
			reader, err := v.GetBody()
			if err != nil {
				return nil, err
			}
			body = reader
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
