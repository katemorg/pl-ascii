package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/katemorg/pl-ascii/ascii"
)

func main() {
	r := mux.NewRouter()
	// health check
	r.HandleFunc("/ping", ascii.Ping)
	r.HandleFunc("/image-to-ascii", ascii.ImageToAscii)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
