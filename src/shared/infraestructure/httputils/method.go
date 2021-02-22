package httputils

import (
	"fmt"
	"net/http"

	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
)

// Verb middleware ensures that petition is made with the rigth http verb
func Verb(acceptedHTTPMethodgo string) middleware.Middleware {
	return func(nextHandler http.HandlerFunc) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			response.Header().Set("Content-Type", "application/json")

			if request.Method != acceptedHTTPMethodgo {
				DispatchNewHttpError(response, fmt.Sprintf("Method %s is not allowed", request.Method), http.StatusMethodNotAllowed)
				return
			}

			nextHandler(response, request)
		}
	}
}
