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

		title := r.FormValue("title")
		if title == "" {
			log.Errorln("missing title from form")
			http.Error(w, "missing title from form", http.StatusBadRequest)
			return
		}
		creator := r.FormValue("creator")
		if creator == "" {
			log.Errorln("missing creator from form")
			http.Error(w, "missing creator from form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("media")
		if err != nil {
			log.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debugf("title: %v\n", title)
		log.Debugf("creator: %v\n", creator)
		log.Debugf("header.Filename: %v\n", header.Filename)
		hash := uuid.NewSHA1(uuid.NameSpaceURL, []byte(header.Filename+creator+title))
		log.Debugf("hash: %v\n", hash.String())
		if err := hls.ConvertStream("hls", hash.String(), file); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := queries.CreateMedia(r.Context(), db.CreateMediaParams{
			Hash:    hash.String(),
			Title:   title,
			Creator: creator,
		}); err != nil {
			log.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type searchResponse struct {
	Count int         `json:"count"`
	Media []db.Medium `json:"media"`
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
