package test

import (
	"goshort"
	"net/http"
	"net/http/httptest"
)

func executeRequest(app goshort.App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}
