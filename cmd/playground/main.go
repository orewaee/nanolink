package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    "127.0.0.1:4000",
		Handler: mux,
	}

	err := server.ListenAndServeTLS("certs/cert.crt", "certs/private.key")
	if err != nil {
		fmt.Println(err)
	}
}
