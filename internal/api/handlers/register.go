package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zSnails/streamx/internal/db"
	"github.com/zSnails/streamx/internal/services/hls"
	"github.com/zSnails/streamx/pkg/logging"
)

var log = logging.Get().WithField("service", "api")

func API(r *mux.Router, queries *db.Queries, conn *pgxpool.Pool) {
	log.Infoln("registering api routes")
	r.HandleFunc("/api/songs", search(queries))
	r.HandleFunc("/api/upload", upload(queries, conn))
	r.HandleFunc("/{hash}/{file-hash}", serve(queries, conn))
}

const (
	_32MB = 33554432
	_1GB  = 1e+9
)

func upload(q *db.Queries, conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tx, err := conn.BeginTx(r.Context(), pgx.TxOptions{})
		if err != nil {
			log.Error(err)
			return
		}
		defer tx.Rollback(r.Context())

		queries := q.WithTx(tx)
		if err := r.ParseMultipartForm(_1GB); err != nil {
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

		log.Debugf("storing media in the database")
		if err := queries.CreateMedia(r.Context(), db.CreateMediaParams{
			Hash:    hash.String(),
			Title:   title,
			Creator: creator,
		}); err != nil {
			log.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hsh := hash.String()
		log.Debugf("converting media to hls data")
		if err := hls.ConvertStream("hls", hsh, file); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dir, err := os.ReadDir(filepath.Join("hls", hsh))
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, file := range dir {
			log.Debugf("file.Name(): %v\n", file.Name())
			content, err := os.ReadFile(filepath.Join("hls", hsh, file.Name()))
			lobs := tx.LargeObjects()
			log.Info("Creating large object")
			objId, err := lobs.Create(r.Context(), 0)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			obj, err := lobs.Open(r.Context(), objId, pgx.LargeObjectModeWrite)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer obj.Close() // BUG: OOOPS FORGOT TO CLOSE THE LARGE OBJECT

			wrote, err := obj.Write(content)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Infof("Wrote %d bytes to large object with the id %d\n", wrote, objId)
			name := fileNameWithoutExtTrimSuffix(file.Name())
			if err := queries.StoreMedia(
				r.Context(),
				db.StoreMediaParams{
					Name: pgtype.Text{
						String: name,
						Valid:  len(name) <= 64,
					},
					Fileoid: pgtype.Uint32{
						Uint32: objId,
						Valid:  true,
					},
				},
			); err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err := os.RemoveAll(filepath.Join("hls", hsh)); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := tx.Commit(r.Context()); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
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

		resp := &searchResponse{}
		if media != nil {
			resp.Media = media
			resp.Count = len(media)
		} else {
			resp.Media = []db.Medium{}
		}
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
