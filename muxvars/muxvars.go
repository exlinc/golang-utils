package muxvars

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetStringVar tries to get the var with the key
func GetStringVar(r *http.Request, key string) string {
	vars := mux.Vars(r)
	if _, ok := vars[key]; !ok {
		return ""
	}
	return vars[key]
}

// GetUnsignedIntegerVar tries to get the var with the key and parse it into a base 10 uint
func GetUnsignedIntegerVar(r *http.Request, key string) uint {
	return uint(GetUnsignedInteger64Var(r, key))
}

// GetUnsignedInteger64Var tries to get the var with the key and parse it into a base 10 uint64
func GetUnsignedInteger64Var(r *http.Request, key string) uint64 {
	vars := mux.Vars(r)
	if _, ok := vars[key]; !ok {
		return 0
	}
	if len(vars[key]) < 1 {
		return 0
	}
	value, err := strconv.ParseUint(vars[key], 10, 64)
	if err != nil {
		return 0
	}
	return uint64(value)
}

// GetIntegerVar tries to get the var with the key and parse it into a base 10 int64
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

// GetURLBase64EncodedStringVar tries to get the var with the key and then parse it into a string using the b64 URL decoder
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
