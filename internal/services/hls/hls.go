package hls

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, fileServer http.Handler) {
	log.Infoln("registering hls service route")
	r.Handle("/playlists/{hash}/{file}", fileServer)
	r.Handle("/videos/{hash}/{file}", fileServer)
}
