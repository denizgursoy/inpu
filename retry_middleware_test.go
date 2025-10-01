package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_RetryMiddleware() {
	httpClient := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(httpClient)

	gock.New(testUrl).
		Get("/").
		Times(2).
		Reply(http.StatusInternalServerError)

	gock.New(testUrl).
		Get("/").
		Times(1).
		Reply(http.StatusOK)

	client := NewWithHttpClient(httpClient).
		UseMiddlewares(RetryMiddleware(2))

	response, err := client.Get(testUrl).Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}
