package posthttpport_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	posthttppost "github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/alejogs4/blog/src/user/domain/user"
	"github.com/gorilla/mux"
)

func preparePostRequest(title, content, tags string, postRespository post.PostRepository) (*httptest.ResponseRecorder, *http.Request, posthttppost.PostControllers) {
	formData := url.Values{}
	formData.Add("title", title)
	formData.Add("content", content)
	formData.Add("tags", tags)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/post", strings.NewReader(formData.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var postCommands application.PostCommands = application.NewPostCommands(postRespository)
	var postQueries application.PostQueries = application.NewPostQueries(postRespository)

	controller := posthttppost.NewPostControllers(postCommands, postQueries)

	rawUser, _ := user.NewUser("id", "Alejandro", "garcia", "alejogs4@gmail.com", "1234567", true)
	userDTO := user.ToDTO(rawUser)
	ctx := context.WithValue(request.Context(), "user", userDTO) //nolint
	ctx = context.WithValue(ctx, "file", "/path/image.jpg")

	return response, request.WithContext(ctx), controller
}

func prepareAddLikeRequest(sentType, postID string, postRespository post.PostRepository) (*httptest.ResponseRecorder, *http.Request, posthttppost.PostControllers) {
	requestBody := []byte(fmt.Sprintf(`{"type": "%v"}`, sentType))

	request := httptest.NewRequest(http.MethodPost, "/api/v1/post/{id}/like", bytes.NewBuffer(requestBody))
	response := httptest.NewRecorder()
	withPostIDRequest := mux.SetURLVars(request, map[string]string{"id": postID})

	request.Header.Set("Content-Type", "application/json")

	rawUser, _ := user.NewUser("id", "Alejandro", "garcia", "alejogs4@gmail.com", "1234567", true)
	userDTO := user.ToDTO(rawUser)
	ctx := context.WithValue(withPostIDRequest.Context(), "user", userDTO) //nolint

	// Here I just noticed that even though I will only use a command I need to pass it a query use case instance, so this could be refactored
	var postCommands application.PostCommands = application.NewPostCommands(postRespository)
	var postQueries application.PostQueries = application.NewPostQueries(postRespository)
	controller := posthttppost.NewPostControllers(postCommands, postQueries)

	return response, request.WithContext(ctx), controller
}

func prepareGetAllPostsRequest(postrepository post.PostRepository) (*httptest.ResponseRecorder, *http.Request, http.HandlerFunc) {
	getAllRequest := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	getAllResponse := httptest.NewRecorder()

	postCommands := application.NewPostCommands(postrepository)
	postQueries := application.NewPostQueries(postrepository)

	postsController := posthttpport.NewPostControllers(postCommands, postQueries)
	getAllPostsRouteController := middleware.Chain(postsController.GetAllPostController, httputils.Verb(http.MethodGet))

	return getAllResponse, getAllRequest, getAllPostsRouteController
}

func existPost(predicate func(post.PostsDTO) bool, posts []post.PostsDTO) bool {
	for _, storedPost := range posts {
		if predicate(storedPost) {
			return true
		}
	}

	return false
}
