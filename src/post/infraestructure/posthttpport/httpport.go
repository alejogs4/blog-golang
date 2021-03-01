package posthttpport

import (
	"net/http"

	"github.com/alejogs4/blog/src/post/application"
	fileupload "github.com/alejogs4/blog/src/post/infraestructure/fileUpload"
	"github.com/alejogs4/blog/src/shared/infraestructure/authentication"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/gorilla/mux"
)

func HandlePostHttpRoutes(router *mux.Router, postCommands application.PostCommands, postQueries application.PostQueries) {
	postController := NewPostControllers(postCommands, postQueries)

	router.HandleFunc("/api/v1/post", middleware.Chain(
		postController.CreatePostController,
		httputils.Verb(http.MethodPost),
		authentication.LoginMiddleare(),
		fileupload.UploadFile("picture", "images"),
	))

	router.HandleFunc("/api/v1/post/{id}/comment", middleware.Chain(
		postController.AddPostComment,
		httputils.Verb(http.MethodPost),
		authentication.LoginMiddleare(),
	))

	router.HandleFunc("/api/v1/comment/{id}", middleware.Chain(
		postController.RemoveComment,
		httputils.Verb(http.MethodDelete),
		authentication.LoginMiddleare(),
	))

	router.HandleFunc("/api/v1/post/{id}/like", middleware.Chain(
		postController.AddPostLikeController,
		httputils.Verb(http.MethodPost),
		authentication.LoginMiddleare(),
	))

	router.HandleFunc("/api/v1/post/{id}", middleware.Chain(
		postController.GetPostByIDController,
		httputils.Verb(http.MethodGet),
	))

	router.HandleFunc("/api/v1/posts", middleware.Chain(
		postController.GetAllPostController,
		httputils.Verb(http.MethodGet),
	))

}
