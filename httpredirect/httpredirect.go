package httpredirect

import "net/http"

// TemporaryRedirect sends a 'StatusTemporaryRedirect' code with the provided URL back to the client. This is useful for non-permanent redirects
func TemporaryRedirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// PermanentRedirect sends a 'StatusPermanentRedirect' code with the provided URL back to the client. This is useful for URLs that have been permanently remapped
func PermanentRedirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}
