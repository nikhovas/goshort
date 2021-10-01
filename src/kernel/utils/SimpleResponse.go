package utils

import (
	"errors"
	"log"
	"net/http"
)

type SimpleResponse struct {
	Status int
	Msg    string
}

func (mr *SimpleResponse) Error() string {
	return mr.Msg
}

func ErrorToResponse(err error, w http.ResponseWriter) {
	var mr *SimpleResponse
	if errors.As(err, &mr) {
		http.Error(w, mr.Msg, mr.Status)
	} else {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
