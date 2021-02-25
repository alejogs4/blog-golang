package posthttpport_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

func TestUnitCreatePostControllerUnit(t *testing.T) {

	t.Run("Should return a bad request code if there are missing field", func(t *testing.T) {
		response, request, controller := preparePostRequest("this is the title", "", "1,2", mockPostRepositoryOK{})
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
		response, request, controller := preparePostRequest("this is the title", "this is the content", "1,2", mockPostRepositoryOK{})

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

func TestUnitAddLikeControllerUnit(t *testing.T) {
	t.Run("Should throw a bad request petition if type property is not sent properly", func(t *testing.T) {
		response, request, controller := prepareAddLikeRequest("invalid-type")
		controller.AddPostLikeController(response, request)

		if response.Code != http.StatusBadRequest {
			t.Errorf("Error: Expected error code %d, received error code %d", http.StatusBadRequest, response.Code)
		}
	})

	t.Run("Should return a created status code if type is sent properly", func(t *testing.T) {
		response, request, controller := prepareAddLikeRequest(like.TLike)
		controller.AddPostLikeController(response, request)

		if response.Code != http.StatusCreated {
			t.Errorf("Error: Expected error code %d, received error code %d", http.StatusCreated, response.Code)
		}
	})
}
