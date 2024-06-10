package main

import (
	"fmt"
	"net/http"
)

type fsAPI struct {
	hits int
}

func (api *fsAPI) mwMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.hits++
		next.ServeHTTP(w, r)
	})
}

func (api *fsAPI) handleDisplayMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", api.hits)))
}

func (api *fsAPI) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	api.hits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
