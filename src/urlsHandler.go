package src

import (
	"github.com/PuerkitoBio/purell"
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix/v3"
	"net/http"
)

func urlsPostHandler(w http.ResponseWriter, r *http.Request) {
	url_, err := UrlFromHttpRequest(w, r)
	if err != nil {
		ErrorToResponse(err, w)
		return
	}

	if url_.Code == 0 {
		url_.Code = 301
	}

	url_.Url, err = purell.NormalizeURLString(url_.Url,
		purell.FlagLowercaseScheme|purell.FlagLowercaseHost|purell.FlagUppercaseEscapes)
	if err != nil {
		ErrorToResponse(&malformedRequest{status: 500, msg: "incorrect URL"}, w)
		return
	}

	if url_.Key == "" {
		url_.Autogenerated = true
		added := false
		for !added {
			potentialKey, err := TryAutogenUrlFromRedis(AppObject.Pool, url_.Url)
			if potentialKey != "" {
				var potentialUrl Url
				if potentialUrl, err = UrlFromRedis(AppObject.Pool, potentialKey); err != nil {
					ErrorToResponse(err, w)
					return
				}

				if potentialUrl.Url == url_.Url {
					url_ = potentialUrl
					added = true
				}
			} else {
				url_.Key = GetNewUrlString(AppObject.Pool)

				if added, err = UrlToRedisInitial(AppObject.Pool, url_); err != nil {
					ErrorToResponse(err, w)
					return
				}
			}
		}
	} else {
		result, err := UrlToRedisInitial(AppObject.Pool, url_)
		if err != nil {
			ErrorToResponse(err, w)
			return
		} else if result {

		}
	}

	UrlToHttpResponse(url_, w)
}

func urlsGetHandler(w http.ResponseWriter, r *http.Request, url Url) {
	UrlToHttpResponse(url, w)
}

func urlsPatchRequest(w http.ResponseWriter, r *http.Request, url Url) {
	newUrl, err := UrlFromHttpRequest(w, r)
	if err != nil {
		ErrorToResponse(err, w)
		return
	}

	url.Url = newUrl.Url
	err = UrlToRedis(AppObject.Pool, url)
	if err != nil {
		ErrorToResponse(err, w)
		return
	}

	UrlToHttpResponse(url, w)
}

func urlsDeleteRequest(w http.ResponseWriter, r *http.Request, key string) {
	err := AppObject.Pool.Do(radix.Cmd(nil, "DEL", key))
	if err != nil {
		ErrorToResponse(err, w)
		return
	}
}

func urlDetailsHandler(w http.ResponseWriter, r *http.Request) {
	url, err := UrlFromRedis(AppObject.Pool, mux.Vars(r)["id"])
	if err != nil {
		ErrorToResponse(err, w)
		return
	}

	if r.Method == "GET" {
		urlsGetHandler(w, r, url)
	} else if r.Method == "PATCH" || r.Method == "PUT" {
		urlsPatchRequest(w, r, url)
	} else { // DELETE
		urlsDeleteRequest(w, r, url.Key)
	}
}

func RegisterUrlsHandlers(router *mux.Router) {
	router.HandleFunc("/urls/", urlsPostHandler).Methods("POST")
	router.HandleFunc("/urls/{id}/", urlDetailsHandler).Methods("GET", "PATCH", "PUT", "DELETE")
}
