package userhttpport

import "github.com/gorilla/mux"

func HandleUserRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/login", loginHandler)
	router.HandleFunc("/api/v1/register", registerHandler)
}
