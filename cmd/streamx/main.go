package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	api "github.com/zSnails/streamx/internal/api/handlers"
	"github.com/zSnails/streamx/internal/api/middleware"
	"github.com/zSnails/streamx/internal/db"
	"github.com/zSnails/streamx/pkg/logging"
)

var (
	log          = logging.Get()
	instancePort string
	dbUser       string
	dbPassword   string
	dbHost       string
)

func getPostgresConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s", dbUser, dbPassword, dbHost)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Warnln(err)
	}

	flag.StringVar(&instancePort, "port", os.Getenv("PORT"), "the port on which to launch this instance")
	flag.StringVar(&dbUser, "postgres-user", os.Getenv("DB_USER"), "the user of the postgres database")
	flag.StringVar(&dbPassword, "postgres-password", os.Getenv("DB_PASSWORD"), "the password of the postgres database")
	flag.StringVar(&dbHost, "postgres-host", os.Getenv("DB_HOST"), "the host of the postgres database")
	flag.Parse()

	if instancePort == "" {
		log.Panicln("the instance port can't be unset either set it on the .env file or pass it as an argument ")
	}

	if dbUser == "" {
		log.Panicln("the user of the postgres database can't be unset either set it on the .env file or pass it as an argument ")
	}

	if dbPassword == "" {
		log.Panicln("the password of the postgres database can't be unset either set it on the .env file or pass it as an argument ")
	}

	if dbHost == "" {
		log.Panicln("the host of the postgres database can't be unset either set it on the .env file or pass it as an argument ")
	}

}

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	mux.CORSMethodMiddleware(r)
	r.Use(middleware.Logger)
	ctx := context.Background()

	log.Infof("connecting to postgres database using user %s\n", dbUser)
	conn, err := pgxpool.New(ctx, getPostgresConnectionString())
	// conn, err := pgx.Connect(ctx, getPostgresConnectionString())
	if err != nil {
		log.Panic(err)
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	queries := db.New(conn)
	api.API(r, queries, conn)
	// r.Handle("/{hash}/{file}", middleware.FilesMW(http.FileServer(http.Dir("hls"))))

	if err := http.ListenAndServe(instancePort, handlers.CORS()(r)); err != nil {
		log.Panic(err)
	}
}
