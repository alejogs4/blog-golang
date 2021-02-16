package userhttpport

import "net/http"

func HandleUserRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/v1/login", loginHandler)
	router.HandleFunc("/api/v1/register", registerHandler)
}
