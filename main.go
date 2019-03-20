package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting...")

	connStr := "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	r := chi.NewRouter()
	r.Use(Health(db))
	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		fmt.Fprintf(w, "Hello World! %s", time.Now())
	})
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {

	})
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	})
	r.Get("/db", func(w http.ResponseWriter, r *http.Request) {
		_, err := db.QueryContext(r.Context(), "SELECT name FROM person")
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			log.Println(err)
		}
		w.Write([]byte("db"))
	})

	s := http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	// if err := db.Close(); err != nil {
	// 	log.Println(err)
	// }
	log.Println("Done DB")
	if err := s.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
	log.Println("Done Server")
}
