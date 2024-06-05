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
    smux.Handle("/", fsrv)
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
