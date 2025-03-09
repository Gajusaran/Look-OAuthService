package main

import (
	"log"
	"net/http"

	"github.com/loginOAuth/router"
)

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}
