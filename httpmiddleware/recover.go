package httpmiddleware

import "net/http"

func RecoverInternalServerError(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
