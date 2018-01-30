package muxvars

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetIntegerVar(r *http.Request, key string) int64 {
	vars := mux.Vars(r)
	if _, ok := vars[key]; !ok {
		return 0
	}
	if len(vars[key]) < 1 {
		return 0
	}
	value, err := strconv.ParseInt(vars[key], 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func GetURLBase64EncodedStringVar(r *http.Request, key string) string {
	vars := mux.Vars(r)
	if _, ok := vars[key]; !ok {
		return ""
	}
	if len(vars[key]) < 1 {
		return ""
	}
	value, err := base64.URLEncoding.DecodeString(vars[key])
	if err != nil {
		return ""
	}
	return string(value)
}
