package test

import (
	"bytes"
	. "gopkg.in/check.v1"
	"goshort"
	"net/http"
)

func (s *MySuite) TestBadToken(c *C) {
	req, err := http.NewRequest("POST", "/urls/", bytes.NewBuffer([]byte("")))
	c.Check(err, Equals, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer ddd")

	response := executeRequest(goshort.AppObject, req)
	c.Check(response.Code, Equals, 401)
}
