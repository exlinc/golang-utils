package hls

import "net/http"

// GetRealmID returns the HLS realm ID used in the request, returning an empty string if the cookie isn't found
func GetRealmID(r *http.Request) string {
	if c, err := r.Cookie("hls_realm_id"); err == nil {
		return c.Value
	}
	return ""
}

// GetHost returns the HLS Host used in the request, returning an empty string if the cookie isn't found
func GetHost(r *http.Request) string {
	if c, err := r.Cookie("hls_host"); err == nil {
		return c.Value
	}
	return ""
}