package inputModules

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"goshort/src/kernel"
	"goshort/src/modules/dbModules"
	"goshort/src/types"
	kernelErrors "goshort/src/types/errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var NotImplementedError = errors.New("not implemented")

func genericKeySupportFalseFunc() bool {
	return false
}

func genericKeySupportTrueFunc() bool {
	return true
}

func createEnvironment() (*dbModules.Generic, *Server) {
	server := &Server{}
	db := &dbModules.Generic{
		GetFunc:               func(_ string) (types.Url, error) { return types.Url{}, NotImplementedError },
		PostFunc:              func(_ types.Url) (types.Url, error) { return types.Url{}, NotImplementedError },
		PatchFunc:             func(_ types.Url) error { return NotImplementedError },
		DeleteFunc:            func(_ string) error { return NotImplementedError },
		GenericKeySupportFunc: genericKeySupportFalseFunc,
		Name:                  "Generic",
	}
	kernelInstance := &kernel.Kernel{}
	kernelInstance.Logger = &kernel.LoggingKernel{Kernel: kernelInstance}
	kernelInstance.Database = &kernel.DatabaseKernel{Kernel: kernelInstance, Database: db}
	kernelInstance.Input = &kernel.InputKernel{Kernel: kernelInstance, Inputs: []types.InputInterface{server}}
	kernelInstance.Middleware = &kernel.MiddlewareKernel{Kernel: kernelInstance}
	kernelInstance.Reconnection = kernel.ReconnectionKernel{Kernel: kernelInstance}
	kernelInstance.Signal = kernel.SignalKernel{Kernel: kernelInstance}

	server.Kernel = kernelInstance
	_ = server.Init(map[string]interface{}{"name": "", "ip": "", "port": 0})

	return db, server
}

func getTestHelper(t *testing.T, getFunc func(key string) (types.Url, error), key string, comparingUrl types.Url,
	comparingReturnCode int) {
	db, server := createEnvironment()
	db.GetFunc = getFunc

	req := httptest.NewRequest(http.MethodGet, "/api/urls/"+key+"/", strings.NewReader(""))
	rec := httptest.NewRecorder()

	server.echo.ServeHTTP(rec, req)
	var resUrl types.Url
	_ = json.NewDecoder(rec.Body).Decode(&resUrl)
	assert.Equal(t, comparingReturnCode, rec.Code)
	assert.Equal(t, comparingUrl, resUrl)
}

func TestSimpleGet(t *testing.T) {
	url := types.Url{Key: "testKey", Url: "http://example.com", Code: 301, Autogenerated: false}
	getFunc := func(key string) (types.Url, error) {
		assert.Equal(t, "testKey", key)
		return url, nil
	}
	getTestHelper(t, getFunc, "testKey", url, http.StatusOK)
}

func TestGetNoExisting(t *testing.T) {
	var url types.Url
	getFunc := func(key string) (types.Url, error) {
		return types.Url{}, kernelErrors.NotFoundError
	}
	getTestHelper(t, getFunc, "", url, http.StatusNotFound)
}

func postTestHelper(t *testing.T, postFunc func(newUrl types.Url) (types.Url, error), genericKeySupport func() bool,
	inputString string, comparingUrl types.Url, comparingReturnCode int) {
	db, server := createEnvironment()
	db.PostFunc = postFunc
	db.GenericKeySupportFunc = genericKeySupport

	req := httptest.NewRequest(http.MethodPost, "/api/urls/", strings.NewReader(inputString))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	server.echo.ServeHTTP(rec, req)
	var resUrl types.Url
	_ = json.NewDecoder(rec.Body).Decode(&resUrl)
	assert.Equal(t, comparingReturnCode, rec.Code)
	assert.Equal(t, comparingUrl, resUrl)
}

func TestSimplePost(t *testing.T) {
	postFunc := func(newUrl types.Url) (types.Url, error) {
		assert.Equal(t, types.Url{Key: "aaa", Url: "https://yandex.ru", Code: 301, Autogenerated: false}, newUrl)
		return newUrl, nil
	}
	postTestHelper(t, postFunc, genericKeySupportFalseFunc, `{"url":"https://yandex.ru","key":"aaa"}`,
		types.Url{Key: "aaa", Url: "https://yandex.ru", Code: 301, Autogenerated: false}, http.StatusCreated)
}

func TestGenericPost(t *testing.T) {
	postFunc := func(newUrl types.Url) (types.Url, error) {
		assert.Equal(t, types.Url{Key: "", Url: "https://yandex.ru", Code: 301, Autogenerated: true}, newUrl)
		newUrl.Key = "a"
		return newUrl, nil
	}
	postTestHelper(t, postFunc, genericKeySupportTrueFunc, `{"url":"https://yandex.ru"}`,
		types.Url{Key: "a", Url: "https://yandex.ru", Code: 301, Autogenerated: true}, http.StatusCreated)
}

func TestAlreadyExistsPost(t *testing.T) {
	postFunc := func(newUrl types.Url) (types.Url, error) {
		return types.Url{}, kernelErrors.AlreadyExistsError
	}
	postTestHelper(t, postFunc, genericKeySupportTrueFunc, "{}", types.Url{}, http.StatusConflict)
}