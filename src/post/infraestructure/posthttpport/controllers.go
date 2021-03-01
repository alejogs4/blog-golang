package posthttpport

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpadapter"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/user/domain/user"
	"github.com/gorilla/mux"
)

type PostControllers struct {
	postCommands application.PostCommands
	postQueries  application.PostQueries
}

func NewPostControllers(postCommands application.PostCommands, postQueries application.PostQueries) PostControllers {
	return PostControllers{postCommands: postCommands, postQueries: postQueries}
}

func (controller PostControllers) CreatePostController(response http.ResponseWriter, request *http.Request) {
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

	userDTO, _ := request.Context().Value("user").(user.UserDTO) //nolint
	userPicture, _ := request.Context().Value("file").(string)

	err := controller.postCommands.CreateNewPost(userDTO.ID, httpBlogPost.Title, httpBlogPost.Content, userPicture, httpBlogPost.Tags)
	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(map[string]string{}, "Post created"), http.StatusCreated)
}

func (controller PostControllers) AddPostComment(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var commentInfo struct {
		Content string `json:"content"`
	}
	postID := mux.Vars(request)["id"]
	userDTO, _ := request.Context().Value("user").(user.UserDTO)

	if err := json.NewDecoder(request.Body).Decode(&commentInfo); err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	commentID, err := controller.postCommands.CreateNewComment(userDTO.ID, postID, commentInfo.Content)
	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	createdComment := map[string]string{"comment_id": commentID}
	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(createdComment, "Comment created"), http.StatusCreated)
}

func (controller PostControllers) AddPostLikeController(response http.ResponseWriter, request *http.Request) {
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

	err = controller.postCommands.AddLike(userDTO.ID, postID, likeInfo.Type)

	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(map[string]string{}, "Ok"), http.StatusCreated)
}

func (controller PostControllers) GetPostByIDController(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	postID := mux.Vars(request)["id"]
	post, err := controller.postQueries.GetPostByID(postID)

	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(post, "Ok"), http.StatusOK)
}

func (controller PostControllers) GetAllPostController(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	posts, err := controller.postQueries.GetAllPosts()
	if err != nil {
		httpError := posthttpadapter.MapPostErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(posts, "Ok"), http.StatusOK)
}
