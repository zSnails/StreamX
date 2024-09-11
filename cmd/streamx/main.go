package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	api "github.com/zSnails/streamx/internal/api/handlers"
	"github.com/zSnails/streamx/internal/api/middleware"
	"github.com/zSnails/streamx/internal/db"
	"github.com/zSnails/streamx/pkg/logging"
)

var log = logging.Get()

func getPostgresConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@localhost", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
}

func main() {
	godotenv.Load()
	r := mux.NewRouter()
	r.NotFoundHandler = middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	mux.CORSMethodMiddleware(r)
	r.Use(middleware.Logger)
	ctx := context.Background()

    log.Infof("connecting to postgres database using user %s\n", os.Getenv("DB_USER"))
	conn, err := pgx.Connect(ctx, getPostgresConnectionString())
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
