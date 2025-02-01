package router

import (
	"github.com/gorilla/mux"
	"github.com/nainanisumit/loginOAuth/handler"
)

func router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register", handler.Register).Methods("GET")
	return router
}
