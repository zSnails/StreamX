package middleware

import "net/http"

func FilesMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: actually check for valid files
		h.ServeHTTP(w, r)
	})
}
