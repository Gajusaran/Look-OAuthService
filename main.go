package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// entry point for project

func main() {
	r := mux.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
