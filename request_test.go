package inpu

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/h2non/gock"
)

type testModel struct {
	Foo string `json:"foo"`
}

var (
	testUrl = "https://x.com"
)

var (
	TestUserName     = "test-user"
	TestUserPassword = "test-password"
)

func (e *ClientSuite) Test_Headers() {
	gock.New(testUrl).
		Get("/").
		MatchHeader(HeaderContentType, MimeTypeJson).
		MatchHeader(HeaderAccept, MimeTypeJson).
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		AcceptJson().
		ContentTypeJson().
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Basic_Authentication() {
	gock.New(testUrl).
		Get("/").
		BasicAuth(TestUserName, TestUserPassword).
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		AuthBasic(TestUserName, TestUserPassword).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Token_Authentication() {
	token := "sdsds"
	gock.New(testUrl).
		Get("/").
		MatchHeader(HeaderAuthorization, "Bearer "+token).
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		AuthToken(token).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Query_Parameters() {
	gock.New(testUrl).
		Get("/").
		MatchParam("is_created", "true").
		MatchParam("foo", "bar").
		MatchParam("float", "1.2").
		MatchParam("float64", "2.2").
		MatchParam("int", "1").
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		QueryBool("is_created", true).
		QueryString("foo", "bar").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

type StarWarsCharacter struct {
	Name      string    `json:"name"`
	Height    string    `json:"height"`
	Mass      string    `json:"mass"`
	HairColor string    `json:"hair_color"`
	SkinColor string    `json:"skin_color"`
	EyeColor  string    `json:"eye_color"`
	BirthYear string    `json:"birth_year"`
	Gender    string    `json:"gender"`
	Homeworld string    `json:"homeworld"`
	Films     []string  `json:"films"`
	Species   []string  `json:"species"`
	Vehicles  []string  `json:"vehicles"`
	Starships []string  `json:"starships"`
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	Url       string    `json:"url"`
}

func (e *ClientSuite) Test_Multiple_Query_Parameterss() {

	response, err := Get("https://swapi.dev/api/people/1").
		QueryInt("foo", 1).
		QueryString("foo1", "bar1").
		Header("foo", "bar").
		Header("foo1", "bar1").
		AuthToken("bar-password").
		Send()

	if response.Status() == http.StatusOK {
		lukeSkywalker := StarWarsCharacter{}
		if err := response.UnmarshalJson(&lukeSkywalker); err != nil {
			log.Fatal(err)
		}
		fmt.Println(lukeSkywalker)
	}

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Multiple_Query_Parameters() {
	// TODO test is wrong
	gock.New(testUrl).
		Get("/").
		MatchParam("foo", "bar1").
		MatchParam("foo", "bar2").
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		QueryString("foo", "bar1").
		QueryString("foo", "bar2").
		Send()

	Get("https://swapi.dev/api/people/1")

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Body_Json_Marshal() {
	gock.New(testUrl).
		Post("/").
		Body(bytes.NewReader([]byte(`{"foo":"bar"}`))).
		Reply(http.StatusOK)

	response, err :=
		Post(testUrl, testModel{Foo: "bar"}).
			Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Body_Reader() {
	gock.New(testUrl).
		Post("/").
		Body(bytes.NewReader([]byte(`{"foo":"bar"}`))).
		Reply(http.StatusOK)

	response, err :=
		Post(testUrl, bytes.NewReader([]byte(`{"foo":"bar"}`))).
			Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}
