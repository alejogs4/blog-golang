package post

import "errors"

// ErrBadPostContent represent when a post is tried to be created with the wrong data
var ErrBadPostContent = errors.New("Content was not provided with the enough data")
var ErrNoFoundPost = errors.New("Post with id was not found")

// Error codes for bad threatment of domain comment rules
var (
	ErrBadCommentContent    = errors.New("Comment was not provided with the enough data")
	ErrInvalidCommentLength = errors.New("Comment cannot be longer than 256 characteres")
	ErrInvalidCommentState  = errors.New("Comment is in removed state")
	ErrUnexistentComment    = errors.New("Referenced comment doesn't exist")
	ErrNoCommentOwner       = errors.New("Only the comment owner can remove it")
)

// ErrInvalidTagInformation is dispatched when a tag is created incorrectly
var ErrInvalidTagInformation = errors.New(("Tag was created incorrectly"))
