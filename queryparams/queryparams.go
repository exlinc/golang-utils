package queryparams

import (
	"net/http"
	"strconv"
)

func GetLimitOffsetQueryParametersSentinel(r *http.Request, sentinel int64) (limit, offset int64) {
	limit, offset = GetLimitOffsetQueryParametersDefaults(r)
	if r.URL.Query().Get("limit") == "" {
		limit = sentinel
	}
	if r.URL.Query().Get("offset") == "" {
		offset = sentinel
	}
	return
}

func GetLimitOffsetQueryParametersDefaults(r *http.Request) (limit, offset int64) {
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit < 0 {
		limit = 10
	}
	return
}

func GetLimitOffsetQueryParametersInt(r *http.Request) (limit, offset int) {
	l, o := GetLimitOffsetQueryParametersDefaults(r)
	return int(l), int(o)
}

func GetLimitOffsetQueryParametersUint(r *http.Request) (limit, offset uint) {
	l, o := GetLimitOffsetQueryParametersDefaults(r)
	return uint(l), uint(o)
}

func GetStringQueryParameter(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
