package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

func newMux(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "User-agent: *\nDisallow: /\n")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if strings.Contains(path, "/") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		id, err := decodeID(path)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		url, err := getURL(db, id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})
	return withLogging(mux)
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", time.Since(start),
			"remote", r.RemoteAddr,
		)
	})
}

func cmdServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", ":8080", "listen address")
	dbPath := fs.String("db", "data/nilink.db", "database path")
	fs.Parse(args)

	db, err := openDB(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("listening", "addr", *addr)
	if err := http.ListenAndServe(*addr, newMux(db)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
