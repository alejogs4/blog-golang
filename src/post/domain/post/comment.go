package post

import "strings"

const (
	// MaxCommentContent how long a comment it's allowed to be
	MaxCommentContent = 256
	// RemovedComment a comment that has been removed
	RemovedComment = "REMOVED"
	// ActiveComment a comment in its by default state
	ActiveComment = "ACTIVE"
)

// Comment .
type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"post_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	State   string `json:"state"` // TODO: this should be a value object
}

// CreateNewComment factory function which create a new comment returning error if data is incorrect
func CreateNewComment(id, postID, userID, content string) (Comment, error) {
	comment := Comment{
		ID:      id,
		PostID:  postID,
		UserID:  userID,
		Content: content,
		State:   ActiveComment,
	}
	error := validateCommentInfo(comment)
	if error != nil {
		return Comment{}, error
	}

	return comment, nil
}

// RemoveComment change comment state to Removed
func (comment *Comment) RemoveComment() error {
	if comment.State == RemovedComment {
		return ErrInvalidCommentState
	}

	comment.State = RemovedComment
	return nil
}

func validateCommentInfo(comment Comment) error {
	commentMetadata := []string{comment.ID, comment.PostID, comment.UserID, comment.Content}
	for _, metadata := range commentMetadata {
		normalizedMetadata := strings.TrimSpace(metadata)
		if normalizedMetadata == "" {
			return ErrBadCommentContent
		}
	}

	if len(comment.Content) > MaxCommentContent {
		return ErrInvalidCommentLength
	}

	if comment.State != RemovedComment && comment.State != ActiveComment {
		return ErrBadCommentContent
	}

	return nil
}
