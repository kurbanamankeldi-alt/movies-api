package errors

import (
	"net/http"
	"log"
)

type HttpError struct {
	Err error
	Message string
	Code int
}

func (e *HttpError) Error() string {
	if e.Err != nil  {
		return e.Err.Error()
	}
	return e.Message
}

type HttpErrorHandler func(http.ResponseWriter, *http.Request) *HttpError

func (fn HttpErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		if e.Err != nil {
			log.Println(e.Err)
		}
		http.Error(w, e.Message, e.Code)
	}
}