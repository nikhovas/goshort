package goshort

import (
	"encoding/json"
	"github.com/mediocregopher/radix/v3"
	"goshort/utils"
	"net/http"
	"strconv"
	"strings"
)

type Url struct {
	Key           string `json:"key"`
	Url           string `json:"url"`
	Code          int    `json:"code"`
	Autogenerated bool   `json:"autogenerated"`
}

func CreateUrlFromRedis(pool *radix.Pool, key string) (Url, error) {
	var preUrl string
	err := pool.Do(radix.Cmd(&preUrl, "GET", key))
	if err != nil {
		return Url{}, err
	}

	if len(preUrl) == 0 {
		return Url{}, err
	}

	s := strings.Split(preUrl, "~")
	code, _ := strconv.Atoi(s[0])

	return Url{
		key,
		strings.Join(s[2:], "~"),
		code,
		s[1] == "1",
	}, nil
}

func CreateUrlFromHttpRequest(w http.ResponseWriter, r *http.Request) (Url, error) {
	var url Url
	err := utils.DecodeJSONBody(w, r, &url)
	return url, err
}

func (url *Url) CreateInRedis(pool *radix.Pool) (bool, error) {
	var p string
	if url.Autogenerated {
		if err := pool.Do(radix.Cmd(&p, "SET", url.Key, strconv.Itoa(url.Code)+"~1~"+url.Url, "NX")); err != nil {
			panic(err)
		}

		if p != "OK" {
			return false, nil
		}

		if err := pool.Do(radix.Cmd(nil, "SET", "$autogen$"+url.Url, url.Key)); err != nil {
			panic(err)
		}

		return true, nil
	} else {
		err := pool.Do(radix.Cmd(&p, "SET", url.Key, strconv.Itoa(url.Code)+"~0~"+url.Url, "NX"))
		if err != nil {
			return false, err
		}
	}
	return p == "OK", nil
}

func (url *Url) UpdateInRedis(pool *radix.Pool) error {
	if url.Autogenerated {
		p := radix.Pipeline(
			radix.Cmd(nil, "SET", url.Key, strconv.Itoa(url.Code)+"~1~"+url.Url),
			radix.Cmd(nil, "SET", "$autogen$"+url.Url, url.Key),
		)

		if err := pool.Do(p); err != nil {
			panic(err)
		}
	} else {
		err := pool.Do(
			radix.Cmd(nil, "SET", url.Key, strconv.Itoa(url.Code)+"~0~"+url.Url))
		if err != nil {
			return err
		}
	}

	return nil
}

func (url *Url) ToHttpResponse(w http.ResponseWriter) {
	js, err := json.Marshal(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func GetAutogenUrlFromRedis(pool *radix.Pool, url string) (string, error) {
	var key string
	if err := pool.Do(radix.Cmd(&key, "GET", "$autogen$"+url)); err != nil {
		return "", err
	}
	return key, nil
}