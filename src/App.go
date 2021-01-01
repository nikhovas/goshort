package src

import (
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix/v3"
	"net/http"
)

type App struct {
	Router *mux.Router
	Pool   *radix.Pool
}

func (a *App) Initialize(network, ip string, poolSize int) {
	var err error
	a.Pool, err = radix.NewPool(network, ip, poolSize)
	if err != nil {
		// handle error
	}

	a.Router = mux.NewRouter().StrictSlash(true)
	RegisterUrlsHandlers(a.Router)
	a.Router.HandleFunc("/{id}", redirector)
}

func (a *App) Run(addr string) {
	http.Handle("/", a.Router)
	_ = http.ListenAndServe(addr, nil)
}
