package userhttpport

import (
	"net/http"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/gorilla/mux"
)

func HandleUserRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/login", middleware.Chain(loginHandler, httputils.Verb(http.MethodPost)))
	router.HandleFunc("/api/v1/register", middleware.Chain(registerHandler, httputils.Verb(http.MethodPost)))
}
