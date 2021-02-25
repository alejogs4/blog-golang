package authentication

import (
	"context"
	"net/http"
	"strings"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/alejogs4/blog/src/shared/infraestructure/token"
)

func LoginMiddleare() middleware.Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			response.Header().Set("Content-Type", "application/json")

			authorization := strings.Fields(request.Header.Get("Authorization"))

			if len(authorization) != 2 {
				httputils.DispatchNewHttpError(response, "Token and authentication type must be present", http.StatusUnauthorized)
				return
			}

			authenticationType := authorization[0]
			if authenticationType != "Bearer" {
				httputils.DispatchNewHttpError(response, "Must be Bearer authentication", http.StatusUnauthorized)
				return
			}

			authenticationToken := authorization[1]
			if authenticationToken == "" {
				httputils.DispatchNewHttpError(response, "Token was not found", http.StatusUnauthorized)
				return
			}

			user, err := token.ValidateToken(authenticationToken)
			if err != nil {
				httputils.DispatchNewHttpError(response, "Token is invalid", http.StatusUnauthorized)
				return
			}

			// Look for what else I can use here
			newContext := context.WithValue(request.Context(), "user", user) // nolint
			f(response, request.WithContext(newContext))
		}
	}
}
