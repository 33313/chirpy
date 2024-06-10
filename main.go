package main

import (
	"fmt"
	"net/http"
)

const (
	ADDRESS = ":8080"
)

func main() {
	api := fsAPI{
		hits: 0,
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))

	mux.Handle("/app/*", api.mwMetrics(http.StripPrefix("/app", fs)))
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("GET /metrics", api.handleDisplayMetrics)
	mux.HandleFunc("/reset", api.handleResetMetrics)
	srv := http.Server{
		Addr:    ADDRESS,
		Handler: mux,
	}

	fmt.Println("Running server on", ADDRESS)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
