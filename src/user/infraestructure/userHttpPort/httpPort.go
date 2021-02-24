package userhttpport

import (
	"net/http"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/alejogs4/blog/src/user/domain/user"
	"github.com/gorilla/mux"
)

func HandleUserRoutes(router *mux.Router, userRepository user.UserRepository) {
	userController := NewUserController(userRepository)

	router.HandleFunc("/api/v1/login", middleware.Chain(userController.LoginHandler, httputils.Verb(http.MethodPost)))
	router.HandleFunc("/api/v1/register", middleware.Chain(userController.RegisterHandler, httputils.Verb(http.MethodPost)))
}
