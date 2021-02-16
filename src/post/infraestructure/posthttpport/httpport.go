package posthttpport

import "github.com/gorilla/mux"

func HandlePostHttpRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/post", createPostController)
	router.HandleFunc("/api/v1/post/{id}", getPostByIDController)
	router.HandleFunc("/api/v1/posts", getAllPostController)
	router.HandleFunc("/api/v1/post/like", addPostLikeController)
}
