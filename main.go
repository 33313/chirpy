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
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /admin/metrics", api.handleDisplayMetrics)
	mux.HandleFunc("/api/reset", api.handleResetMetrics)
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
