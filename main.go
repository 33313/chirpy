package main

import (
	"fmt"
	"github.com/myshkovsky/chirpy/internal/database"
	"net/http"
)

const (
	ADDRESS = ":8080"
)

func main() {
	newDB, err := database.NewDB("./database.json")
	if err != nil {
		panic(err)
	}
	api := fsAPI{
		hits: 0,
		db:   newDB,
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))

	mux.Handle("/app/*", api.mwMetrics(http.StripPrefix("/app", fs)))
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /admin/metrics", api.handleDisplayMetrics)
	mux.HandleFunc("/api/reset", api.handleResetMetrics)

	mux.HandleFunc("POST /api/chirps", api.handlePostChirp)
	mux.HandleFunc("GET /api/chirps", api.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", api.handleGetChirp)

	mux.HandleFunc("GET /api/users/{id}", api.handleGetUser)
	mux.HandleFunc("POST /api/users", api.handlePostUser)

	srv := http.Server{
		Addr:    ADDRESS,
		Handler: mux,
	}

	fmt.Println("Running server on", ADDRESS)
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
