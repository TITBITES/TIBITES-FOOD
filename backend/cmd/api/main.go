package main

import (
	"log"
	"net/http"

	"local.dev/foodapp/internal/httpapi"
)

func main() {
	r := httpapi.NewRouter()
	addr := ":8080"
	log.Printf("starting http server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
