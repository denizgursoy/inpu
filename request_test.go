package inpu

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type testModel struct {
	Foo string `json:"foo" xml:"foo"`
}

var (
	TestUserName     = "test-user"
	TestUserPassword = "test-password"
	testData         = testModel{Foo: "bar"}
	testDataAsJson   = `{"foo":"bar"}`
	testDataAsXml    = `<testModel><foo>bar</foo></testModel>`
)

func (c *ClientSuite) Test_Headers() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal(MimeTypeJson, request.Header.Get(HeaderAccept))
		c.Require().Equal(MimeTypeJson, request.Header.Get(HeaderContentType))
		c.Require().Equal("test-user", request.Header.Get(HeaderUserAgent))
	}))
	defer server.Close()

	err := Get(server.URL).
		AcceptJson().
		ContentTypeJson().
		UserAgent("test-user").
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Basic_Authentication() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("Basic dGVzdC11c2VyOnRlc3QtcGFzc3dvcmQ=", request.Header.Get(HeaderAuthorization))
	}))
	defer server.Close()

	err := Get(server.URL).
		AuthBasic(TestUserName, TestUserPassword).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Token_Authentication() {
	token := "sdsds"
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("Bearer "+token, request.Header.Get(HeaderAuthorization))
	}))
	defer server.Close()

	err := Get(server.URL).
		AuthToken(token).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Query_Parameters() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/?float=1.2000000476837158&float64=2.2&foo=bar+test+encoded&int=1&is_created=true", request.RequestURI)
	}))
	defer server.Close()

	err := Get(server.URL).
		QueryBool("is_created", true).
		QueryString("foo", "bar test encoded").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Multiple_Query_Parameters() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/?foo=bar1&foo=bar2", request.RequestURI)
	}))
	defer server.Close()

	err := Get(server.URL).
		QueryString("foo", "bar1").
		QueryString("foo", "bar2").
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Body_Json_Marshal() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(request.Body)
		c.Require().NoError(err)
		c.Require().Equal(testDataAsJson, string(all))
	}))
	defer server.Close()

	err := Post(server.URL, BodyJson(testData)).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Request_Timeout() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := Post(server.URL, BodyJson(testData)).
		TimeOutIn(100*time.Millisecond).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().ErrorIs(err, context.DeadlineExceeded)
}

func (c *ClientSuite) Test_Body_Reader() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(request.Body)
		c.Require().NoError(err)
		c.Require().Equal(testDataAsJson, string(all))
	}))
	defer server.Close()

	err := Post(server.URL, BodyReader(bytes.NewReader([]byte(testDataAsJson)))).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Multiple_Chose_Correct_Reply_Behaviour() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	expectedError := errors.New("correct reply was executed")

	c.T().Run("should select StatusIs over StatusIsSuccess", func(t *testing.T) {
		err := Get(server.URL).
			QueryString("foo", "bar1").
			QueryString("foo", "bar2").
			OnReplyIf(StatusIsSuccess, ThenReturnError(errors.New("unexpected status"))).
			OnReplyIf(StatusIs(http.StatusOK), ThenReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
	c.T().Run("should select StatusIs over StatusIsOneOf", func(t *testing.T) {
		err := Get(server.URL).
			QueryString("foo", "bar1").
			QueryString("foo", "bar2").
			OnReplyIf(StatusIsOneOf(http.StatusOK, http.StatusAccepted), ThenReturnError(errors.New("unexpected status"))).
			OnReplyIf(StatusIs(http.StatusOK), ThenReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
	c.T().Run("should select StatusIsOneOf over StatusIsSuccess", func(t *testing.T) {
		err := Get(server.URL).
			QueryString("foo", "bar1").
			QueryString("foo", "bar2").
			OnReplyIf(StatusIsSuccess, ThenReturnError(errors.New("unexpected status"))).
			OnReplyIf(StatusIsOneOf(http.StatusOK, http.StatusAccepted), ThenReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
	c.T().Run("should select StatusIsOneOf over StatusAny ", func(t *testing.T) {
		err := Get(server.URL).
			QueryString("foo", "bar1").
			QueryString("foo", "bar2").
			OnReplyIf(StatusAny, ThenReturnError(errors.New("unexpected status"))).
			OnReplyIf(StatusIsOneOf(http.StatusOK, http.StatusAccepted), ThenReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
}

// Non-pointer value tests
func (c *ClientSuite) Test_QueryInt8() {
	c.T().Parallel()
	req := Get("").QueryInt8("foo", 8)
	c.Require().Equal([]string{"8"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryInt16() {
	c.T().Parallel()
	req := Get("").QueryInt16("foo", 16)
	c.Require().Equal([]string{"16"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryInt32() {
	c.T().Parallel()
	req := Get("").QueryInt32("foo", 32)
	c.Require().Equal([]string{"32"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryInt() {
	c.T().Parallel()
	req := Get("").QueryInt("foo", 42)
	c.Require().Equal([]string{"42"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryInt64() {
	c.T().Parallel()
	req := Get("").QueryInt64("foo", 999)
	c.Require().Equal([]string{"999"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint8() {
	c.T().Parallel()
	req := Get("").QueryUint8("foo", 255)
	c.Require().Equal([]string{"255"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint16() {
	c.T().Parallel()
	req := Get("").QueryUint16("foo", 65535)
	c.Require().Equal([]string{"65535"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint32() {
	c.T().Parallel()
	req := Get("").QueryUint32("foo", 4294967295)
	c.Require().Equal([]string{"4294967295"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint() {
	c.T().Parallel()
	req := Get("").QueryUint("foo", 100)
	c.Require().Equal([]string{"100"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint64() {
	c.T().Parallel()
	req := Get("").QueryUint64("foo", 999)
	c.Require().Equal([]string{"999"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryFloat32() {
	c.T().Parallel()
	req := Get("").QueryFloat32("foo", 3.14)
	c.Require().Equal([]string{"3.140000104904175"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryFloat64() {
	c.T().Parallel()
	req := Get("").QueryFloat64("foo", 3.14159)
	expected := strconv.FormatFloat(3.14159, 'f', -1, 64)
	c.Require().Equal([]string{expected}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryBool_True() {
	c.T().Parallel()
	req := Get("").QueryBool("foo", true)
	c.Require().Equal([]string{"true"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryBool_False() {
	c.T().Parallel()
	req := Get("").QueryBool("foo", false)
	c.Require().Equal([]string{"false"}, req.queries["foo"])
}

func (c *ClientSuite) Test_Query() {
	c.T().Parallel()
	req := Get("").QueryString("foo", "bar")
	c.Require().Equal([]string{"bar"}, req.queries["foo"])
}

// Pointer value tests - non-nil
func (c *ClientSuite) Test_QueryInt8Ptr_NonNil() {
	c.T().Parallel()
	val := int8(8)
	req := Get("").QueryInt8Ptr("foo", &val)
	c.Require().Equal([]string{"8"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryInt16Ptr_NonNil() {
	c.T().Parallel()
	val := int16(16)
	req := Get("").QueryInt16Ptr("foo", &val)
	c.Require().Equal([]string{"16"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryInt32Ptr_NonNil() {
	c.T().Parallel()
	val := int32(32)
	req := Get("").QueryInt32Ptr("foo", &val)
	c.Require().Equal([]string{"32"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryIntPtr_NonNil() {
	c.T().Parallel()
	val := 42
	req := Get("").QueryIntPtr("foo", &val)
	c.Require().Equal([]string{"42"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryInt64Ptr_NonNil() {
	c.T().Parallel()
	val := int64(999)
	req := Get("").QueryInt64Ptr("foo", &val)
	c.Require().Equal([]string{"999"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint8Ptr_NonNil() {
	c.T().Parallel()
	val := uint8(255)
	req := Get("").QueryUint8Ptr("foo", &val)
	c.Require().Equal([]string{"255"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint16Ptr_NonNil() {
	c.T().Parallel()
	val := uint16(65535)
	req := Get("").QueryUint16Ptr("foo", &val)
	c.Require().Equal([]string{"65535"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint32Ptr_NonNil() {
	c.T().Parallel()
	val := uint32(4294967295)
	req := Get("").QueryUint32Ptr("foo", &val)
	c.Require().Equal([]string{"4294967295"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUintPtr_NonNil() {
	c.T().Parallel()
	val := uint(100)
	req := Get("").QueryUintPtr("foo", &val)
	c.Require().Equal([]string{"100"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryUint64Ptr_NonNil() {
	c.T().Parallel()
	val := uint64(999)
	req := Get("").QueryUint64Ptr("foo", &val)
	c.Require().Equal([]string{"999"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryFloat32Ptr_NonNil() {
	c.T().Parallel()
	val := float32(3.14)
	req := Get("").QueryFloat32Ptr("foo", &val)
	c.Require().Equal([]string{"3.140000104904175"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryFloat64Ptr_NonNil() {
	c.T().Parallel()
	val := 3.14159
	req := Get("").QueryFloat64Ptr("foo", &val)
	expected := strconv.FormatFloat(3.14159, 'f', -1, 64)
	c.Require().Equal([]string{expected}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryBoolPtr_NonNil_True() {
	c.T().Parallel()
	val := true
	req := Get("").QueryBoolPtr("foo", &val)
	c.Require().Equal([]string{"true"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryBoolPtr_NonNil_False() {
	c.T().Parallel()
	val := false
	req := Get("").QueryBoolPtr("foo", &val)
	c.Require().Equal([]string{"false"}, req.queries["foo"])
}

func (c *ClientSuite) Test_QueryStringPtr_NonNil() {
	c.T().Parallel()
	val := "bar"
	req := Get("").QueryStringPtr("foo", &val)
	c.Require().Equal([]string{"bar"}, req.queries["foo"])
}

// Pointer value tests - nil (should NOT add anything to queries)
func (c *ClientSuite) Test_QueryInt8Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryInt8Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryInt16Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryInt16Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryInt32Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryInt32Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryIntPtr_Nil() {
	c.T().Parallel()
	req := Get("").QueryIntPtr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryInt64Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryInt64Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryUint8Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryUint8Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryUint16Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryUint16Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryUint32Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryUint32Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryUintPtr_Nil() {
	c.T().Parallel()
	req := Get("").QueryUintPtr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryUint64Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryUint64Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryFloat32Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryFloat32Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryFloat64Ptr_Nil() {
	c.T().Parallel()
	req := Get("").QueryFloat64Ptr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryBoolPtr_Nil() {
	c.T().Parallel()
	req := Get("").QueryBoolPtr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

func (c *ClientSuite) Test_QueryStringPtr_Nil() {
	c.T().Parallel()
	req := Get("").QueryStringPtr("foo", nil)
	c.Require().Len(req.queries, 0)
	c.Require().Empty(req.queries.Get("foo"))
}

// Test method chaining
func (c *ClientSuite) Test_QueryMethodChaining() {
	c.T().Parallel()
	req := Get("").
		QueryInt("id", 123).
		QueryString("name", "test").
		QueryBool("active", true)

	c.Require().Equal([]string{"123"}, req.queries["id"])
	c.Require().Equal([]string{"test"}, req.queries["name"])
	c.Require().Equal([]string{"true"}, req.queries["active"])
}

// Test that nil pointers don't interfere with subsequent calls
func (c *ClientSuite) Test_QueryPtrNil_ThenNonNil() {
	c.T().Parallel()
	val := int64(555)
	req := Get("").
		QueryInt64Ptr("foo", nil).
		QueryInt64Ptr("bar", &val)

	// foo should not exist
	c.Require().Empty(req.queries.Get("foo"))
	// bar should exist
	c.Require().Equal([]string{"555"}, req.queries["bar"])
}

func (c *ClientSuite) Test_Multiple_Response_Handler() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	expectedError := errors.New("expected error")
	result := testModel{}
	err := Post(server.URL, nil).
		OnReplyIf(StatusIs(http.StatusOK),
			ThenUnmarshalJsonTo(&result),
			ThenReturnDefaultError,
			ThenReturnError(expectedError),
		).Send()
	expectedResponse := testModel{
		Foo: "bar",
	}
	defaultError := &DefaultError{}

	c.Require().Equal(expectedResponse, result)
	c.Require().Error(err)
	c.Require().ErrorIs(err, expectedError)
	c.Require().ErrorAs(err, &defaultError)
}
