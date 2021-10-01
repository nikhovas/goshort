package inputModules

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"goshort/kernel"
	"goshort/types"
	kernelErrors "goshort/types/errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	types.ModuleBase
	echo         *echo.Echo
	ip           string
	port         string
	moduleName   string
	Kernel       *kernel.Kernel
	nativeClosed bool
}

func CreateServer(kernel *kernel.Kernel) types.InputInterface {
	return &Server{Kernel: kernel}
}

func (server *Server) urlsPostHandler(c echo.Context) error {
	var newUrl types.Url
	if err := json.NewDecoder(c.Request().Body).Decode(&newUrl); err != nil {
		return nil
	}

	postedUrl, err := server.Kernel.Database.Post(newUrl)
	// add exists error support
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, postedUrl)
}

func (server *Server) urlsPatchRequest(c echo.Context) error {
	id := c.Param("id")
	url_, err := server.Kernel.Database.Get(id)
	if err != nil {
		return err
	}

	var newUrl types.Url
	if err := json.NewDecoder(c.Request().Body).Decode(&newUrl); err != nil {
		return err
	}

	url_.Url = newUrl.Url
	if err := server.Kernel.Database.Patch(url_); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, url_)
}

func (server *Server) urlsGetHandler(c echo.Context) error {
	id := c.Param("id")
	url_, err := server.Kernel.Database.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, url_)
}

func (server *Server) urlsDeleteRequest(c echo.Context) error {
	id := c.Param("id")
	err := server.Kernel.Database.Delete(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (server *Server) registerUrlsHandlers(g *echo.Group) {
	g.POST("/urls/", server.urlsPostHandler)
	g.GET("/urls/:id/", server.urlsGetHandler)
	g.PATCH("/urls/:id/", server.urlsPatchRequest)
	g.PUT("/urls/:id/", server.urlsPatchRequest)
	g.DELETE("/urls/:id/", server.urlsDeleteRequest)
}

func (server *Server) redirect(c echo.Context) error {
	id := c.Param("id")
	urlVal, _ := server.Kernel.Database.Get(id)
	if urlVal.Url == "" {
		return c.HTML(http.StatusNotFound, "<h1>Not found.</h1>")
	}
	return c.Redirect(urlVal.Code, urlVal.Url)
}

func (server *Server) mainPage(c echo.Context) error {
	return c.HTML(http.StatusNotFound, "<h1>Main page</h1>")
}

type ServerLog struct {
	ClientIp string
	Status   int
	Endpoint string
	Method   string
	Type     string
	Error    types.Log
}

func (l *ServerLog) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["name"] = "Input.Server.ConnectionLog"
	data["type"] = l.Type
	data["endpoint"] = l.Endpoint
	data["clientIp"] = l.ClientIp
	data["method"] = l.Method
	data["status"] = l.Status
	if l.Error != nil {
		data["error"] = l.Error
	}
	return data
}

func (server *Server) mainLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err types.Log
		err = nil
		if err2 := next(c); err2 != nil {
			//httpErr, ok := err2.(echo.HTTPError)
			err = &kernelErrors.SimpleErrorWrapper{Err: err2}
		}

		le := &ServerLog{
			ClientIp: c.RealIP(),
			Status:   c.Response().Status,
			Endpoint: c.Request().RequestURI,
			Method:   c.Request().Method,
			Type:     "",
			Error:    err,
		}
		if err == nil {
			le.Type = "log"
		} else {
			le.Type = "error"
		}
		_ = server.Kernel.Logger.Send(le)
		return nil
	}
}

func (server *Server) Init(config map[string]interface{}) error {
	if err := server.ModuleBase.Init(config); err != nil {
		return err
	}

	server.echo = echo.New()
	server.moduleName = config["name"].(string)
	server.ip = config["ip"].(string)
	server.port = strconv.Itoa(config["port"].(int))
	server.echo.HideBanner = true
	server.echo.HidePort = true

	api := server.echo.Group("api")
	api.POST("/urls/", server.urlsPostHandler)
	api.GET("/urls/:id/", server.urlsGetHandler)
	api.PATCH("/urls/:id/", server.urlsPatchRequest)
	api.PUT("/urls/:id/", server.urlsPatchRequest)
	api.DELETE("/urls/:id/", server.urlsDeleteRequest)

	server.echo.GET("/:id", server.redirect)
	server.echo.GET("/", server.mainPage)
	server.echo.Use(server.mainLoggingMiddleware)
	return nil
}

func (server *Server) Run(wg *sync.WaitGroup) error {
	go func() {
		defer wg.Done()
		server.Kernel.SetModuleRunState(server)
		defer server.Kernel.SetModuleStopState(server)
		err := server.echo.Start(server.ip + ":" + server.port)
		if err != nil {
			if !(err == http.ErrServerClosed && server.nativeClosed) {
				_ = server.Kernel.Logger.SendError(err)
			}
		}
	}()
	return nil
}

func (server *Server) Stop() error {
	server.nativeClosed = true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := server.echo.Shutdown(ctx)
	return err
}

func (server *Server) GetName() string {
	return server.moduleName
}

func (server *Server) GetType() string {
	return "Input.Server"
}
