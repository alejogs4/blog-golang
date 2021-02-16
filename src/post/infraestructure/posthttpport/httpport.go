package posthttpport

import (
	"github.com/alejogs4/blog/src/shared/infraestructure/authentication"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/gorilla/mux"
)

func HandlePostHttpRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/post", middleware.Chain(createPostController, authentication.LoginMiddleare()))
	router.HandleFunc("/api/v1/post/{id}", getPostByIDController)
	router.HandleFunc("/api/v1/posts", getAllPostController)
	router.HandleFunc("/api/v1/post/{id}/like", middleware.Chain(addPostLikeController, authentication.LoginMiddleare()))
}
