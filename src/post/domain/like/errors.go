package like

import "errors"

var (
	// ErrBadLikeContent is dispatched when a like is tried to be created with not enough data
	ErrBadLikeContent = errors.New("Like was not provided with the enough data")
	// ErrInvalidLikeType is dispatched when a like type is different to Dislike or Like
	ErrInvalidLikeType = errors.New("Like must be either Dislike or Like")
	// ErrInvalidLikeState is dispatched when like state is different to active or removed
	ErrInvalidLikeState = errors.New("Like must be either active or removed")
)
