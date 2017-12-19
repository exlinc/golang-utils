package htmlhttp

import (
	"fmt"
	"net/http"
)

// InternalServerErrorView returns a simple HTML page with the internal server error status code
func InternalServerErrorView(w http.ResponseWriter, r *http.Request) {
	HTMLWriter(w, r, "<div>Internal server error</div>", http.StatusInternalServerError)
}

// UnauthorizedErrorView returns a simple HTML page with the unauthorized status code
func UnauthorizedErrorView(w http.ResponseWriter, r *http.Request) {
	HTMLWriter(w, r, "<div>Unauthorized request error</div>", http.StatusUnauthorized)
}

// HTMLWriter provides a wrapper function to send HTML-string responses back over an http.ResponseWriter
func HTMLWriter(w http.ResponseWriter, r *http.Request, html string, status int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", html)
}