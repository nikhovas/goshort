package src

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	url, _ := UrlFromRedis(AppObject.Pool, mux.Vars(r)["id"])
	http.Redirect(w, r, url.Url, url.Code)
}
