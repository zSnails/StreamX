package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func FilesMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: actually check for valid files
		vars := mux.Vars(r)
		log.Debugf("vars: %v\n", vars)
		h.ServeHTTP(w, r)
	})
}
