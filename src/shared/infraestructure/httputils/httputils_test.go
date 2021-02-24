package httputils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
)

func TestVerbMiddlewareUnit(t *testing.T) {
	bareHandler := func(response http.ResponseWriter, request *http.Request) {
		httputils.DispatchNewResponse(response, map[string]string{}, http.StatusOK)
	}

	handlerWithVerb := middleware.Chain(bareHandler, httputils.Verb(http.MethodGet))
	testCases := []struct {
		Name       string
		Verb       string
		StatusCode int
	}{
		{Name: "Should return method not allowed if verb is different than expected", Verb: http.MethodPost, StatusCode: http.StatusMethodNotAllowed},
		{Name: "Should return ok code if verb is the expected", Verb: http.MethodGet, StatusCode: http.StatusOK},
	} // Table test, I'll try to refactor other tests cases to use this pattern

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(c.Verb, "/any-url", nil)

			handlerWithVerb(response, request)

			if response.Code != c.StatusCode {
				t.Errorf("Error: expected status code %d, received %d", c.StatusCode, response.Code)
			}
		})
	}
}
