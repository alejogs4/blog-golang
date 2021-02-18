package post

import "errors"

// Error codes for violation of post entity rules
var (
	ErrBadPostContent     = errors.New("Post: Content was not provided with the enough data")
	ErrNoFoundPost        = errors.New("Post: Post with id was not found")
	ErrMissingPostPicture = errors.New("Post: Post picture is mandatory")
)

// Error codes for bad threatment of domain comment rules
var (
	ErrBadCommentContent    = errors.New("Comment: Comment was not provided with the enough data")
	ErrInvalidCommentLength = errors.New("Comment: Comment cannot be longer than 256 characteres")
	ErrInvalidCommentState  = errors.New("Comment: Comment is in removed state")
	ErrUnexistentComment    = errors.New("Comment: Referenced comment doesn't exist")
	ErrNoCommentOwner       = errors.New("Comment: Only the comment owner can remove it")
)

// ErrInvalidTagInformation is dispatched when a tag is created incorrectly
var ErrInvalidTagInformation = errors.New(("Tag was created incorrectly"))
