package httpcookies

import (
	"fmt"
	"net/http"
	"time"
)

// SetJavascriptAccessibleCookie should NOT UNDER ANY CIRCUMSTANCES be used for sensitive credentials! It sets an HTTP cookie that *is* accessible through the client Javascript. This should only be done for insecure variables that are useful to the client
func SetJavascriptAccessibleCookie(w http.ResponseWriter, expires time.Time, domain string, cookieName string, cookieValue interface{}) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   fmt.Sprint(cookieValue),
		Domain: domain,
		Path:    "/",
		Expires: expires,
	})
}

// SetHTTPOnlyCookie sets a cookie that is only visible to the browser, as opposed to the client JS itself. This makes these cookies safe for access tokens and other credentials
func SetHTTPOnlyCookie(w http.ResponseWriter, expires time.Time, domain string, cookieName string, cookieValue interface{}) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    fmt.Sprint(cookieValue),
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
		Domain: domain,
	})
}

// DeleteCookie removes a cookie from the clients browser, useful for things like getting rid of a client's session cookie
func DeleteCookie(w http.ResponseWriter, domain string, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
		Path:   "/",
		Domain: domain,
	})
}
