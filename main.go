package main

import (
	"fmt"
	"net/http"
)

const (
    ADDRESS = ":8080"
)

func main() {
    smux := http.NewServeMux()
    fsrv := http.FileServer(http.Dir("."))
    smux.Handle("/app/*", http.StripPrefix("/app", fsrv))
    smux.HandleFunc("/healthz", handleHealthz)
    srv := http.Server{
        Addr: ADDRESS,
        Handler: smux,
    }
    fmt.Println("Running server on", ADDRESS)
    err := srv.ListenAndServe()
    if err != nil {
        panic(err)
    }
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}
