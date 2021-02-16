package posthttpport

import "net/http"

func HandlePostHttpRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/v1/post", createPostController)
	router.HandleFunc("/api/v1/posts", getAllPostController)
}
