package main

import (
	"log"
	"net/http"

	"github.com/Gajusaran/Look-OAuthService/router"
)

// entry point for project

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}
