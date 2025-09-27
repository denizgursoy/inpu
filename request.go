package inpu

import (
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Req struct {
	method          string
	rawUrl          string
	headers         map[string]string
	body            any
	timeOutDuration time.Duration
	followRedirect  bool
}

func Get(url string) *Req {
	return &Req{
		method: http.MethodGet,
		rawUrl: url,
	}
}

func Post(url string, body any) *Req {
	return &Req{
		method: http.MethodPost,
		rawUrl: url,
		body:   body,
	}
}

func Delete(url string, body any) *Req {
	return &Req{
		method: http.MethodDelete,
		rawUrl: url,
		body:   body,
	}
}

func Put(url string, body any) *Req {
	return &Req{
		method: http.MethodPut,
		rawUrl: url,
		body:   body,
	}
}

func Patch(url string, body any) *Req {
	return &Req{
		method: http.MethodPatch,
		rawUrl: url,
		body:   body,
	}
}

func (r *Req) UseHttp11() *Req {
	// &http.Transport{ ForceAttemptHTTP2: false, // disable HTTP/2 }
	return r
}

// UseHttpClient can be used in the testing
func (r *Req) UseHttpClient(client *http.Client) *Req {
	// &http.Transport{ ForceAttemptHTTP2: false, // disable HTTP/2 }
	return r
}

func (r *Req) UseTransport(transport *http.Transport) *Req {
	// &http.Transport{ ForceAttemptHTTP2: false, // disable HTTP/2 }
	return r
}

func (r *Req) UseTlsConfig(tlsConfig *tls.Config) *Req {
	// &http.Transport{ ForceAttemptHTTP2: false, // disable HTTP/2 }
	return r
}

func (r *Req) InsecureSkipVerify() *Req {
	// tlsConfig := &tls.Config{InsecureSkipVerify: true}
}

func (r *Req) Header(key, val string) *Req {
	r.addHeader(key, val)
	return r
}

func (r *Req) ContentTypeJson() *Req {
	r.ContentType("application/json")

	return r
}

func (r *Req) ContentTypeText() *Req {
	r.ContentType("text/plain")

	return r
}

func (r *Req) ContentTypeHtml() *Req {
	r.ContentType("text/html")

	return r
}

func (r *Req) ContentType(contentType string) *Req {
	r.addHeader("Content-Type", contentType)

	return r
}

func (r *Req) FollowRedirect() *Req {
	r.followRedirect = true
	return r
}

func (r *Req) AuthBasic(username, password string) *Req {
	cred := username + ":" + password
	r.addHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(cred)))

	return r
}

func (r *Req) AuthToken(token string) *Req {
	r.addHeader("Authorization", "Bearer "+token)
	return r
}

func (r *Req) AcceptJson() *Req {
	r.addHeader("Accepts", "application/json")
	return r
}

func (r *Req) TimeOutIn(duration time.Duration) *Req {
	r.timeOutDuration = duration
	return r
}

func (r *Req) addQueryValue(key, value string) *Req {
	u, err := url.Parse(r.rawUrl)
	if err != nil {
		// If the URL is invalid, return as-is
		return r
	}
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	r.rawUrl = u.String()
	return r
}

func (r *Req) addHeader(key, value string) *Req {
	r.headers[key] = value

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

func (r *Req) Send() (Response, error) {
	return Response{}, nil
}
