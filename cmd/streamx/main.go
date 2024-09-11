package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zSnails/streamx/internal/api/middleware"
	"github.com/zSnails/streamx/internal/services/hls"
	"github.com/zSnails/streamx/pkg/logging"
)

var log = logging.Get()

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	mux.CORSMethodMiddleware(r)
	log.Infoln(os.Getwd())

	r.Use(middleware.Logger)
	// TODO: remove this piece of shit, only used for development
	r.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		hls.Convert("hls", "furina", "furina-bg-processed.mp4")
	})
	r.Handle("/{hash}/{file}", middleware.FilesMW(http.FileServer(http.Dir("hls"))))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Panic(err)
	}
}
