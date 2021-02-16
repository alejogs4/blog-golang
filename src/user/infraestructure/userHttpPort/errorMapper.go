package userhttpport

import (
	"errors"
	"net/http"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/user/domain/user"
)

func mapUserErrorToHttpError(err error) httputils.HttpError {
	if errors.Is(err, user.ErrBadUserData) {
		return httputils.HttpError{Message: err.Error(), Status: http.StatusBadRequest}
	}

	if errors.Is(err, user.ErrTooShortUserPassword) {
		return httputils.HttpError{Message: err.Error(), Status: http.StatusBadRequest}
	}

	if errors.Is(err, user.ErrInvalidUserLogin) {
		return httputils.HttpError{Message: err.Error(), Status: http.StatusBadRequest}
	}

	return httputils.HttpError{Message: "Something went wrong", Status: http.StatusInternalServerError}
}
