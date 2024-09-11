package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zSnails/streamx/internal/db"
	"github.com/zSnails/streamx/internal/services/hls"
	"github.com/zSnails/streamx/pkg/logging"
)

var log = logging.Get().WithField("service", "api")

func API(r *mux.Router, queries *db.Queries) {
	log.Infoln("registering api routes")
	r.HandleFunc("/api/songs", search(queries))
	r.HandleFunc("/api/upload", upload(queries))
}

const _32MB = 33554432

func upload(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(_32MB); err != nil {
			log.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, header, err := r.FormFile("media")
		if err != nil {
			log.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debugf("header.Filename: %v\n", header.Filename)
        hash := uuid.NewSHA1(uuid.NameSpaceURL, []byte(header.Filename))
		log.Debugf("hash: %v\n", hash.String())
		hls.ConvertStream("hls", hash.String(), file)

		// queries.CreateMedia(r.Context(), db.CreateMediaParams{
		// 	Hash:    "",
		// 	Title:   "",
		// 	Creator: "",
		// })
	}
}

func search(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("q")
		var (
			media []db.Medium
			err   error
		)

		if term == "" {
			media, err = queries.AllMedia(r.Context())
		} else {
			media, err = queries.FindMedia(r.Context(), term)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(media); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
