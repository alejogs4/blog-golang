package posthttpport

import (
	"encoding/json"
	"net/http"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpadapter"
	"github.com/alejogs4/blog/src/post/infraestructure/postrepository"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
)

var postCommands application.PostCommands = application.NewPostCommands(postrepository.PostgresRepository{})
var postQueries application.PostQueries = application.NewPostQueries(postrepository.PostgresRepository{})

func createPostController(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var httpBlogPost post.Post

	err := json.NewDecoder(request.Body).Decode(&httpBlogPost)

	if err != nil {
		httputils.DispatchNewHttpError(response, "All fields must be sent", http.StatusBadRequest)
		return
	}

	err = postCommands.CreateNewPost(httpBlogPost.UserID, httpBlogPost.Title, httpBlogPost.Content, httpBlogPost.Picture, httpBlogPost.Tags)
	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	responseContent, _ := json.Marshal(map[string]string{"Message": "Post created"})

	response.WriteHeader(http.StatusCreated)
	response.Write(responseContent)
}

func getAllPostController(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	posts, err := postQueries.GetAllPosts()
	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(posts, "Ok"), http.StatusOK)
}
