package router

import (
	"github.com/Gajusaran/Look-OAuthService/handler"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/verifyotp", handler.VerifyOTP).Methods("POST")
	router.HandleFunc("/resendotp", handler.ResendOTP).Methods("POST")
	router.HandleFunc("/login", handler.ResendOTP).Methods("POST")
	router.HandleFunc("/refresh-token", handler.RefreshTokenHandler).Methods("POST")
	return router
}
