package test

import (
	"UrlShortener/src"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"UrlShortener/utils"
	. "gopkg.in/check.v1"
)

func executeRequest(app src.App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	//App src.App
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {

	utils.SetupViper()
	src.AppObject = src.App{}
	src.AppObject.Initialize("tcp", "127.0.0.1:6379", 10)
	//s.App.Initialize("tcp", "127.0.0.1", 10)
}

func (s *MySuite) TestHelloWorld(c *C) {
	var jsonStr = []byte(`{
	"key": "c",
	"url": "https://yandex.ru"
}`)
	req, _ := http.NewRequest("POST", "/urls/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(src.AppObject, req)
	c.Check(http.StatusCreated, Equals, response.Code)

	//checkResponseCode(t, http.StatusCreated, response.Code)
	//
	//locClient := client.LocationClient{Host: s.host}
	//c.Check(42, Equals, 42)
}
