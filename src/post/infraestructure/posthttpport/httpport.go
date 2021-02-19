package posthttpport

import (
	"github.com/alejogs4/blog/src/post/application"
	fileupload "github.com/alejogs4/blog/src/post/infraestructure/fileUpload"
	"github.com/alejogs4/blog/src/shared/infraestructure/authentication"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/gorilla/mux"
)

func HandlePostHttpRoutes(router *mux.Router, postCommands application.PostCommands, postQueries application.PostQueries) {
	postController := NewPostControllers(postCommands, postQueries)

	router.HandleFunc("/api/v1/post", middleware.Chain(
		postController.CreatePostController,
		authentication.LoginMiddleare(),
		fileupload.UploadFile("picture", "images"),
	))
	router.HandleFunc("/api/v1/post/{id}", postController.GetPostByIDController)
	router.HandleFunc("/api/v1/posts", postController.GetAllPostController)
	router.HandleFunc("/api/v1/post/{id}/like", middleware.Chain(postController.AddPostLikeController, authentication.LoginMiddleare()))
}
