package src

import (
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix/v3"
	"log"
	"net/http"
)

func faviconHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

type App struct {
	Router *mux.Router
	Pool   *radix.Pool
}

func (a *App) Initialize(network, ip string, poolSize int) {
	var err error
	a.Pool, err = radix.NewPool(network, ip, poolSize)
	if err != nil {
		log.Panicln("Can't connect to redis database. Aborting.")
	} else {
		log.Println("Connected to redis database.")
	}

	a.Router = mux.NewRouter().StrictSlash(true)
	a.Router.HandleFunc("/favicon.ico", faviconHandler)
	RegisterUrlsHandlers(a.Router)
	a.Router.HandleFunc("/{id}", Redirect)

}

func (a *App) Run(addr string) {
	http.Handle("/", a.Router)
	_ = http.ListenAndServe(addr, nil)
}
