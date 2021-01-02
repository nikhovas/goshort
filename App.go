package goshort

import (
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix/v3"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	Pool   *radix.Pool
}

var AppObject = App{}

func Redirect(w http.ResponseWriter, r *http.Request) {
	url, _ := CreateUrlFromRedis(AppObject.Pool, mux.Vars(r)["id"])
	http.Redirect(w, r, url.Url, url.Code)
}

func faviconHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (a *App) Initialize() {
	ip := viper.GetString("redis.ip")
	poolSize := viper.GetInt("redis.poolSize")
	var err error
	a.Pool, err = radix.NewPool("tcp", ip, poolSize)
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

func (a *App) Run() {
	http.Handle("/", a.Router)
	_ = http.ListenAndServe(":"+viper.GetString("port"), nil)
}
