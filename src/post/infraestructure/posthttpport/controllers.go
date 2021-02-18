package posthttpport

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpadapter"
	"github.com/alejogs4/blog/src/post/infraestructure/postrepository"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/user/domain/user"
	"github.com/gorilla/mux"
)

var postCommands application.PostCommands = application.NewPostCommands(postrepository.PostgresRepository{})
var postQueries application.PostQueries = application.NewPostQueries(postrepository.PostgresRepository{})

func createPostController(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var httpBlogPost post.Post
	httpBlogPost.Content = request.FormValue("content")
	httpBlogPost.Title = request.FormValue("title")

	tags := strings.Split(strings.TrimSpace(request.FormValue("tags")), ",")
	postTags := make([]post.Tag, len(tags))

	for _, tag := range tags {
		postTags = append(postTags, post.Tag{ID: tag, Content: tag})
	}
	httpBlogPost.Tags = postTags

	userDTO, _ := request.Context().Value("user").(user.UserDTO)
	userPicture, _ := request.Context().Value("file").(string)

	err := postCommands.CreateNewPost(userDTO.ID, httpBlogPost.Title, httpBlogPost.Content, userPicture, httpBlogPost.Tags)
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
		Type string `json:"type"`
	}

	postID := mux.Vars(request)["id"]
	// Improvement check for errors here
	userDTO, _ := request.Context().Value("user").(user.UserDTO)

	err := json.NewDecoder(request.Body).Decode(&likeInfo)
	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	err = postCommands.AddLike(userDTO.ID, postID, likeInfo.Type)

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
