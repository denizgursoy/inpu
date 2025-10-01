package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_RequestModifierMiddleware() {
	httpClient := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(httpClient)

	gock.New(testUrl).
		Get("/").
		MatchHeader("foo", "bar").
		MatchHeader("foo3", "bar3").
		MatchParam("foo1", "bar1").
		MatchParam("foo2", "bar2").
		Reply(http.StatusOK)

	client := NewWithHttpClient(httpClient).
		UseMiddlewares(RequestModifierMiddleware(func(request *http.Request) {
			request.Header.Add("foo", "bar")
			query := request.URL.Query()
			query.Add("foo1", "bar1")
			request.URL.RawQuery = query.Encode()
		}))

	response, err := client.
		Get(testUrl).
		QueryString("foo2", "bar2").
		Header("foo3", "bar3").
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())

}
