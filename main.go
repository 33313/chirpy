package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/myshkovsky/chirpy/internal/database"
)

const (
	ADDRESS = ":8080"
)

func main() {
	godotenv.Load()
	newDB := database.NewDB("./database.json")
	api := API{
		hits:      0,
		db:        newDB,
		jwtSecret: os.Getenv("JWT_SECRET"),
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
	mux.HandleFunc("DELETE /api/chirps/{id}", api.handleDeleteChirp)

	mux.HandleFunc("POST /api/users", api.handlePostUser)
	mux.HandleFunc("GET /api/users/{id}", api.handleGetUser)
	mux.HandleFunc("PUT /api/users", api.handlePutUser)

	mux.HandleFunc("POST /api/login", api.handleLogin)
	mux.HandleFunc("POST /api/refresh", api.handleRefresh)
	mux.HandleFunc("POST /api/revoke", api.handleRevoke)

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
