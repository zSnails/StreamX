package api

import (
	"bufio"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zSnails/streamx/internal/db"
)

func serve(q *db.Queries, conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx, err := conn.Begin(r.Context())
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())
		queries := q.WithTx(tx)
		vars := mux.Vars(r)
		lobs := tx.LargeObjects()

		log.Debugf("vars: %+v\n", vars)

		hash, ok := vars["file-hash"]
		if !ok {
			log.Error("missing file hash")
			http.Error(w, "missing file hash", http.StatusInternalServerError)
			return
		}
		log.Debugf("Using hash %s\n", hash)
		hash = fileNameWithoutExtTrimSuffix(hash)
		media, err := queries.GetStoredMedia(r.Context(), pgtype.Text{
			String: hash,
			Valid:  true,
		})
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lob, err := lobs.Open(r.Context(), media.Fileoid.Uint32, pgx.LargeObjectModeRead)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Debugf("lob: %+v\n", lob)
		buffered := bufio.NewReader(lob)
		wrote, err := buffered.WriteTo(w)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Infof("Wrote %d bytes\n", wrote)

		if err := tx.Commit(r.Context()); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
