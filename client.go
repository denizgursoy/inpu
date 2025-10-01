package inpu

import (
	"context"
	"encoding/base64"
	"maps"
	"net/http"
	netUrl "net/url"
	"strconv"
	"time"
)

type Client struct {
	headers         http.Header
	queries         netUrl.Values
	userClient      *http.Client
	timeOutDuration time.Duration
}

func New() *Client {
	return &Client{
		headers:    make(http.Header),
		queries:    make(netUrl.Values),
		userClient: &http.Client{},
	}
}

func NewWithHttpClient(client *http.Client) *Client {
	return &Client{
		headers:    make(http.Header),
		queries:    make(netUrl.Values),
		userClient: client,
	}
}

func (c *Client) Get(url string) *Req {
	return getReq(context.Background(), url, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) GetCtx(ctx context.Context, url string) *Req {
	return getReq(ctx, url, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) Post(url string, body any) *Req {
	return postReq(context.Background(), url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) PostCtx(ctx context.Context, url string, body any) *Req {
	return postReq(ctx, url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) Delete(url string, body any) *Req {
	return deleteReq(context.Background(), url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) DeleteCtx(ctx context.Context, url string, body any) *Req {
	return deleteReq(ctx, url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) Put(url string, body any) *Req {
	return putReq(context.Background(), url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) PutCtx(ctx context.Context, url string, body any) *Req {
	return putReq(ctx, url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) Patch(url string, body any) *Req {
	return patchReq(context.Background(), url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}
func (c *Client) PatchCtx(ctx context.Context, url string, body any) *Req {
	return patchReq(ctx, url, body, maps.Clone(c.headers), maps.Clone(c.queries), c.userClient, c.timeOutDuration)
}

func (c *Client) Header(key, val string) *Client {
	c.addHeader(key, val)
	return c
}

func (c *Client) UseMiddlewares(mws ...Middleware) *Client {
	c.setDefaultTransportIfEmpty()

	for i := range mws {
		middleware := mws[i]
		if middleware != nil {
			c.userClient.Transport = middleware(c.userClient.Transport)
		}
	}

	return c
}

func (c *Client) setDefaultTransportIfEmpty() *Client {
	if c.userClient.Transport == nil {
		c.userClient.Transport = http.DefaultTransport
	}

	return c
}

func (c *Client) ContentTypeJson() *Client {
	c.ContentType(MimeTypeJson)

	return c
}

func (c *Client) ContentTypeText() *Client {
	c.ContentType(MimeTypeText)

	return c
}

func (c *Client) ContentTypeHtml() *Client {
	c.ContentType(MimeTypeHtml)

	return c
}

func (c *Client) ContentType(contentType string) *Client {
	c.addHeader(HeaderContentType, contentType)

	return c
}

func (c *Client) AuthBasic(username, password string) *Client {
	cred := username + ":" + password
	c.addHeader(HeaderAuthorization, "Basic "+base64.StdEncoding.EncodeToString([]byte(cred)))

	return c
}

func (c *Client) AuthToken(token string) *Client {
	c.addHeader(HeaderAuthorization, "Bearer "+token)
	return c
}

func (c *Client) AcceptJson() *Client {
	c.addHeader(HeaderAccept, MimeTypeJson)
	return c
}

func (c *Client) TimeOutIn(duration time.Duration) *Client {
	c.timeOutDuration = duration

	return c
}

func (c *Client) addQueryValue(key, value string) *Client {
	c.queries.Add(key, value)

	return c
}

func (c *Client) addHeader(key, value string) *Client {
	c.headers.Add(key, value)

	return c
}

func (c *Client) QueryInt8(name string, v int8) *Client {
	return c.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (c *Client) QueryInt16(name string, v int16) *Client {
	return c.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (c *Client) QueryInt32(name string, v int32) *Client {
	return c.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (c *Client) QueryInt(name string, v int) *Client {
	return c.addQueryValue(name, strconv.FormatInt(int64(v), 10))
}

func (c *Client) QueryInt64(name string, v int64) *Client {
	return c.addQueryValue(name, strconv.FormatInt(v, 10))
}

func (c *Client) QueryUint8(name string, v uint8) *Client {
	return c.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (c *Client) QueryUint16(name string, v uint16) *Client {
	return c.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (c *Client) QueryUint32(name string, v uint32) *Client {
	return c.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (c *Client) QueryUint(name string, v uint) *Client {
	return c.addQueryValue(name, strconv.FormatUint(uint64(v), 10))
}

func (c *Client) QueryUint64(name string, v uint64) *Client {
	return c.addQueryValue(name, strconv.FormatUint(v, 10))
}

func (c *Client) QueryFloat32(name string, v float32) *Client {
	return c.addQueryValue(name, strconv.FormatFloat(float64(v), 'f', -1, 64))
}

func (c *Client) QueryFloat64(name string, v float64) *Client {
	return c.addQueryValue(name, strconv.FormatFloat(v, 'f', -1, 64))
}

func (c *Client) QueryBool(name string, v bool) *Client {
	return c.addQueryValue(name, strconv.FormatBool(v))
}

func (c *Client) QueryString(name string, v string) *Client {
	return c.addQueryValue(name, v)
}

func (c *Client) QueryInt8Ptr(name string, v *int8) *Client {
	if v == nil {
		return c
	}
	return c.QueryInt8(name, *v)
}

func (c *Client) QueryInt16Ptr(name string, v *int16) *Client {
	if v == nil {
		return c
	}
	return c.QueryInt16(name, *v)
}

func (c *Client) QueryInt32Ptr(name string, v *int32) *Client {
	if v == nil {
		return c
	}
	return c.QueryInt32(name, *v)
}

func (c *Client) QueryIntPtr(name string, v *int) *Client {
	if v == nil {
		return c
	}
	return c.QueryInt(name, *v)
}

func (c *Client) QueryInt64Ptr(name string, v *int64) *Client {
	if v == nil {
		return c
	}
	return c.QueryInt64(name, *v)
}

func (c *Client) QueryUint8Ptr(name string, v *uint8) *Client {
	if v == nil {
		return c
	}
	return c.QueryUint8(name, *v)
}

func (c *Client) QueryUint16Ptr(name string, v *uint16) *Client {
	if v == nil {
		return c
	}
	return c.QueryUint16(name, *v)
}

func (c *Client) QueryUint32Ptr(name string, v *uint32) *Client {
	if v == nil {
		return c
	}
	return c.QueryUint32(name, *v)
}

func (c *Client) QueryUintPtr(name string, v *uint) *Client {
	if v == nil {
		return c
	}
	return c.QueryUint(name, *v)
}

func (c *Client) QueryUint64Ptr(name string, v *uint64) *Client {
	if v == nil {
		return c
	}
	return c.QueryUint64(name, *v)
}

func (c *Client) QueryFloat32Ptr(name string, v *float32) *Client {
	if v == nil {
		return c
	}
	return c.QueryFloat32(name, *v)
}

func (c *Client) QueryFloat64Ptr(name string, v *float64) *Client {
	if v == nil {
		return c
	}
	return c.QueryFloat64(name, *v)
}

func (c *Client) QueryBoolPtr(name string, v *bool) *Client {
	if v == nil {
		return c
	}
	return c.QueryBool(name, *v)
}

func (c *Client) QueryStringPtr(name string, v *string) *Client {
	if v == nil {
		return c
	}
	return c.QueryString(name, *v)
}
