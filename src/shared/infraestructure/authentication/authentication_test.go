package authentication_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejogs4/blog/src/shared/infraestructure/authentication"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
)

func TestLoginMiddlewareUnit(t *testing.T) {
	handler := func(response http.ResponseWriter, request *http.Request) {

	}
	protectedHandler := middleware.Chain(handler, authentication.LoginMiddleare())

	testCases := []struct {
		Name       string
		TokenValue string
		StatusCode int
		Message    string
	}{
		{Name: "Should be present auth type and token", TokenValue: "", StatusCode: http.StatusUnauthorized, Message: "Token and authentication type must be present"},
		{Name: "Should be present a valid auth type: Bearer", TokenValue: "notvalid  any", StatusCode: http.StatusUnauthorized, Message: "Must be Bearer authentication"},
		// This one should be seen carefully, see how generate a valid token and test the validity of a random one
		{Name: "Should be present an invalid token", TokenValue: "Bearer eytoken", StatusCode: http.StatusUnauthorized, Message: "Token is invalid"},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			var errorResponse struct {
				Message string `json:"message"`
			}
			request := httptest.NewRequest(http.MethodGet, "/any-url", nil)
			response := httptest.NewRecorder()
			request.Header.Set("Authorization", c.TokenValue)

			protectedHandler(response, request)
			if response.Code != c.StatusCode {
				t.Errorf("Error: expected status code %d, received %d", c.StatusCode, response.Code)
			}

			json.NewDecoder(response.Body).Decode(&errorResponse)
			if errorResponse.Message != c.Message {
				t.Errorf("Error: expected message %v, received %v", c.Message, errorResponse.Message)
			}
		})
	}
}
