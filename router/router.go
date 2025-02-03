package router

import (
	"github.com/Gajusaran/Look-OAuthService/handler"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/verifyotp", handler.VerifyOTP).Methods("POST")
	return router
}
