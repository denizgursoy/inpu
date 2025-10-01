package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_RequestModifierMiddleware() {
	gock.New(testUrl).
		Get("/").
		MatchParam("foo", "bar").
		Reply(http.StatusOK)

	client := New().
		UseMiddlewares(RequestModifierMiddleware(func(request *http.Request) {
			request.Header.Add("foo", "bar")
		}))

	response, err := client.Get(testUrl).Send()
	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())

}
