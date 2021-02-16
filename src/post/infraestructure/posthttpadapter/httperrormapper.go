package posthttpadapter

import (
	"errors"
	"net/http"

	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
)

func MapPostErrorToHttpError(err error) httputils.HttpError {
	if errors.Is(err, post.ErrBadPostContent) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: post.ErrBadPostContent.Error()}
	}

	if errors.Is(err, post.ErrBadCommentContent) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: post.ErrBadCommentContent.Error()}
	}

	if errors.Is(err, post.ErrInvalidCommentLength) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: post.ErrInvalidCommentLength.Error()}
	}

	if errors.Is(err, post.ErrInvalidCommentState) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: post.ErrInvalidCommentState.Error()}
	}

	if errors.Is(err, post.ErrUnexistentComment) {
		return httputils.HttpError{Status: http.StatusNotFound, Message: post.ErrUnexistentComment.Error()}
	}

	if errors.Is(err, post.ErrNoCommentOwner) {
		return httputils.HttpError{Status: http.StatusForbidden, Message: post.ErrNoCommentOwner.Error()}
	}

	if errors.Is(err, post.ErrInvalidTagInformation) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: post.ErrInvalidTagInformation.Error()}
	}

	if errors.Is(err, like.ErrBadLikeContent) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: like.ErrBadLikeContent.Error()}
	}

	if errors.Is(err, like.ErrInvalidLikeType) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: like.ErrInvalidLikeType.Error()}
	}

	if errors.Is(err, like.ErrInvalidLikeState) {
		return httputils.HttpError{Status: http.StatusBadRequest, Message: like.ErrInvalidLikeState.Error()}
	}

	return httputils.HttpError{Status: http.StatusInternalServerError, Message: "Something went wrong"}
}
