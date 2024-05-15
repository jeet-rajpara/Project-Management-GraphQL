package database

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (db *sql.DB, err error) {

	// connect to database
	err_in_loadenv := godotenv.Load()
	if err_in_loadenv != nil {
		log.Fatal("Error loading .env file")
	}
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

// Middleware Store CockroachDB connection object in to Request Context

func Middleware(db *sql.DB) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = context.WithValue(r.Context(), "db", db)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return
}

// func Middleware(db *sql.DB, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var ctx = context.WithValue(r.Context(), "db", db)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
