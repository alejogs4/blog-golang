package fileupload

import (
	"errors"
	"net/http"

	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
)

func mapFileErrorToHttpError(err error) httputils.HttpError {
	if errors.Is(err, post.ErrMissingPostPicture) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: post.ErrMissingPostPicture.Error()}
	}

	if errors.Is(err, errCopyingFile) {
		return httputils.HttpError{Status: http.StatusInternalServerError, Message: errCopyingFile.Error()}
	}

	return httputils.HttpError{Status: http.StatusInternalServerError, Message: "Something went wrong"}
}
