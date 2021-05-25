package inputModules

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"goshort/kernel"
	kernelErrors "goshort/kernel/errors"
	"net/http"
	"strconv"
)

type CantDecodeRequestError struct{}

func (e *CantDecodeRequestError) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["name"] = "Input.Server.CantDecodeRequest"
	data["type"] = "error"
	return data
}

func (e *CantDecodeRequestError) Error() string {
	return "Error Input.Server.CantDecodeRequest"
}

type CantEncodeRequestError struct{}

func (e *CantEncodeRequestError) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["name"] = "Input.Server.CantEncodeRequest"
	data["type"] = "error"
	return data
}

func (e *CantEncodeRequestError) Error() string {
	return "Error Input.Server.CantEncodeRequest"
}

func CantEncodeRequestErrorWrapper(err error) *CantEncodeRequestError {
	if err != nil {
		return &CantEncodeRequestError{}
	}
	return nil
}

type Server struct {
	echo       *echo.Echo
	ip         string
	port       string
	moduleName string
	Kernel     *kernel.Kernel
}

func (server *Server) urlsPostHandler(c echo.Context) error {
	var newUrl kernel.Url
	if err := json.NewDecoder(c.Request().Body).Decode(&newUrl); err != nil {
		return &CantDecodeRequestError{}
	}

	postedUrl, err := server.Kernel.Post(newUrl)
	// add exists error support
	if err != nil {
		return err
	}

	return CantEncodeRequestErrorWrapper(c.JSON(http.StatusOK, postedUrl))
}

func (server *Server) urlsPatchRequest(c echo.Context) error {
	id := c.Param("id")
	url_, err := kernel.Object.Get(id)
	if err != nil {
		return err
	}

	var newUrl kernel.Url
	if err := json.NewDecoder(c.Request().Body).Decode(&newUrl); err != nil {
		return err
	}

	url_.Url = newUrl.Url
	if err := server.Kernel.Patch(url_); err != nil {
		return err
	}

	return CantEncodeRequestErrorWrapper(c.JSON(http.StatusOK, url_))
}

func (server *Server) urlsGetHandler(c echo.Context) error {
	id := c.Param("id")
	url_, err := server.Kernel.Get(id)
	if err != nil {
		return err
	}

	return CantEncodeRequestErrorWrapper(c.JSON(http.StatusOK, url_))
}

func (server *Server) urlsDeleteRequest(c echo.Context) error {
	id := c.Param("id")
	url_, err := server.Kernel.Get(id)
	if err != nil {
		return err
	}

	err = server.Kernel.Delete(url_)
	if err != nil {
		return err
	}

	return CantEncodeRequestErrorWrapper(c.NoContent(http.StatusOK))
}

func (server *Server) registerUrlsHandlers(g *echo.Group) {
	g.POST("/urls/", server.urlsPostHandler)
	g.GET("/urls/:id/", server.urlsGetHandler)
	g.PATCH("/urls/:id/", server.urlsPatchRequest)
	g.PUT("/urls/:id/", server.urlsPatchRequest)
	g.DELETE("/urls/:id/", server.urlsDeleteRequest)
}

//func CheckTokenMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		correctToken := viper.GetString("token")
//		if correctToken != "" {
//			token := r.Header.Get("Authorization")
//			if token == "" {
//				utils.ErrorToResponse(&utils.SimpleResponse{Status: http.StatusUnauthorized, Msg: "Need auth credentials"}, w)
//				return
//			} else if "Bearer "+correctToken != token {
//				utils.ErrorToResponse(&utils.SimpleResponse{Status: http.StatusUnauthorized, Msg: "Bad auth credentials"}, w)
//				return
//			}
//		}
//		next(w, r)
//	}
//}

func (server *Server) redirect(c echo.Context) error {
	id := c.Param("id")
	urlVal, _ := server.Kernel.Get(id)
	if urlVal.Url == "" {
		return c.HTML(http.StatusNotFound, "<h1>Not found.</h1>")
	}
	return c.Redirect(urlVal.Code, urlVal.Url)
}

func (server *Server) mainPage(c echo.Context) error {
	return c.HTML(http.StatusNotFound, "<h1>Main page</h1>")
}

//func faviconHandler(w http.ResponseWriter, _ *http.Request) {
//	w.WriteHeader(http.StatusNotFound)
//}

func errorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}
		switch err {
		case kernelErrors.UrlNotFoundError:
			return c.NoContent(http.StatusNotFound)
		default:
			return err
		}
	}
}

type ServerLog struct {
	ClientIp string
	Status   int
	Endpoint string
	Method   string
	Type     string
	Error    kernel.AdvancedError
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
		var err kernel.AdvancedError
		err = nil
		if err2 := next(c); err2 != nil {
			httpErr, ok := err2.(echo.HTTPError)
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
		server.Kernel.Log(le)
		return nil
	}
}

func (server *Server) Init(config map[string]interface{}) error {
	server.echo = echo.New()
	server.moduleName = config["name"].(string)
	server.ip = config["ip"].(string)
	server.port = strconv.Itoa(config["port"].(int))
	server.echo.HideBanner = true

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

func (server *Server) Run() error {
	viper.GetString("port")
	return server.echo.Start(server.ip + ":" + server.port)
}

func (server *Server) GetName() string {
	return server.moduleName
}

func (server *Server) GetType() string {
	return "Server"
}
