package posthttpport

import (
	"encoding/json"
	"net/http"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpadapter"
	"github.com/alejogs4/blog/src/post/infraestructure/postrepository"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/gorilla/mux"
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

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(map[string]string{}, "Post created"), http.StatusCreated)
}

func addPostLikeController(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var likeInfo struct {
		UserID string `json:"user_id"`
		Type   string `json:"type"`
	}
	postID := mux.Vars(request)["id"]

	err := json.NewDecoder(request.Body).Decode(&likeInfo)
	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	err = postCommands.AddLike(likeInfo.UserID, postID, likeInfo.Type)

	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(map[string]string{}, "Ok"), http.StatusCreated)
}

func getPostByIDController(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	postID := mux.Vars(request)["id"]
	post, err := postQueries.GetPostByID(postID)

	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(post, "Ok"), http.StatusOK)
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
