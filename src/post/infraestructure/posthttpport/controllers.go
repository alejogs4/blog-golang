package posthttpport

import (
	"encoding/json"
	"io/ioutil"
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

	// This could be refactored
	file, _, err := request.FormFile("picture")
	if err != nil {
		httputils.DispatchNewHttpError(response, "Something went wrong reading the picture", http.StatusBadRequest)
		return
	}
	defer file.Close()

	newFile, err := ioutil.TempFile("images", "upload-*.jpeg")
	if err != nil {
		httputils.DispatchNewHttpError(response, "Something went wrong copying the picture", http.StatusInternalServerError)
		return
	}
	picturePath := "/" + newFile.Name()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		httputils.DispatchNewHttpError(response, "Something went wrong copying the picture", http.StatusInternalServerError)
		return
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		httputils.DispatchNewHttpError(response, "Something went wrong copying the picture", http.StatusInternalServerError)
		return
	}
	//

	userDTO, _ := request.Context().Value("user").(user.UserDTO)
	err = postCommands.CreateNewPost(userDTO.ID, httpBlogPost.Title, httpBlogPost.Content, picturePath, httpBlogPost.Tags)
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
