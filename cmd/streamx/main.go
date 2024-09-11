package main

import (
	"context"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	api "github.com/zSnails/streamx/internal/api/handlers"
	"github.com/zSnails/streamx/internal/api/middleware"
	"github.com/zSnails/streamx/internal/db"
	"github.com/zSnails/streamx/pkg/logging"
)

var log = logging.Get()

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	mux.CORSMethodMiddleware(r)
	r.Use(middleware.Logger)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgresql://postgres@localhost")
	if err != nil {
		log.Panic(err)
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	queries := db.New(conn)
	api.API(r, queries)
	r.Handle("/{hash}/{file}", middleware.FilesMW(http.FileServer(http.Dir("hls"))))

	if err := http.ListenAndServe(":8080", handlers.CORS()(r)); err != nil {
		log.Panic(err)
	}
}
