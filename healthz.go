package main

import (
	"database/sql"
	"net/http"
	"strings"
)

func Health(db *sql.DB) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.EqualFold(r.URL.Path, "/healthz") {
				err := db.PingContext(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"ok"}`))
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
