package main

import (
	"net/http"
)

func httpErrCode(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func httpErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
