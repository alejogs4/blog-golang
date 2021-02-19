package posthttpport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	posthttppost "github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/user/domain/user"
)

func TestUnitCreatePostController(t *testing.T) {

	t.Run("Should return a bad request code if there are missing field", func(t *testing.T) {
		response, request, controller := preparePostRequest("this is the title", "", "1,2")
		controller.CreatePostController(response, request)

		if response.Code != http.StatusBadRequest {
			t.Errorf("Error: Expected status code: %d, received status code: %d", http.StatusBadRequest, response.Code)
		}

		var responseJSON struct {
			Message string `json:"message"`
		}
		json.NewDecoder(response.Body).Decode(&responseJSON)
		if responseJSON.Message != post.ErrBadPostContent.Error() {
			t.Errorf("Error: expected error message %v, received error message %v", post.ErrBadPostContent.Error(), responseJSON.Message)
		}
	})

	t.Run("Should return StatusCreated if information was correctly provided", func(t *testing.T) {
		response, request, controller := preparePostRequest("this is the title", "this is the content", "1,2")

		controller.CreatePostController(response, request)
		if response.Code != http.StatusCreated {
			t.Errorf("Error: Expected status code: %d, received status code: %d", http.StatusCreated, response.Code)
		}

		var receivedResponse struct {
			Data    interface{} `json:"data"`
			Message string      `json:"message"`
		}
		json.NewDecoder(response.Body).Decode(&receivedResponse)

		expectedMessage := "Post created"
		if receivedResponse.Message != expectedMessage {
			t.Errorf("Error: Expected message %v, received message %v", expectedMessage, receivedResponse.Message)
		}
	})
}

func preparePostRequest(title, content, tags string) (*httptest.ResponseRecorder, *http.Request, posthttppost.PostControllers) {
	formData := url.Values{}
	formData.Add("title", title)
	formData.Add("content", content)
	formData.Add("tags", tags)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/post", strings.NewReader(formData.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockRepository := mockPostRepositoryOK{}
	var postCommands application.PostCommands = application.NewPostCommands(mockRepository)
	var postQueries application.PostQueries = application.NewPostQueries(mockRepository)

	controller := posthttppost.NewPostControllers(postCommands, postQueries)

	rawUser, _ := user.NewUser("id", "Alejandro", "garcia", "alejogs4@gmail.com", "1234567", true)
	userDTO := user.ToDTO(rawUser)
	ctx := context.WithValue(request.Context(), "user", userDTO)
	ctx = context.WithValue(ctx, "file", "/path/image.jpg")

	return response, request.WithContext(ctx), controller
}
