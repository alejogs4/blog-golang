package post

import (
	"strings"

	"github.com/alejogs4/blog/src/post/domain/like"
)

type Post struct {
	ID       string      `json:"id"`
	UserID   string      `json:"user_id"`
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	Picture  string      `json:"picture"`
	Tags     []Tag       `json:"tags"`
	Comments []Comment   `json:"comments"`
	Likes    []like.Like `json:"likes"`
}

// CreateNewPost will verify that right data was provided and return a new instance of the post if so
func CreateNewPost(id, userID, title, content, picture string, comments []Comment, tags []Tag, likes []like.Like) (Post, error) {
	normalizedID := strings.TrimSpace(id)
	normalizedUserID := strings.TrimSpace(userID)
	normalizedTitle := strings.TrimSpace(title)
	normalizedContent := strings.TrimSpace(content)
	normalizedPicture := strings.TrimSpace(picture)

	if normalizedID == "" || normalizedUserID == "" || normalizedTitle == "" || normalizedContent == "" || normalizedPicture == "" {
		return Post{}, ErrBadPostContent
	}

	return Post{
		ID:       id,
		UserID:   userID,
		Title:    title,
		Content:  content,
		Picture:  picture,
		Tags:     tags,
		Comments: comments,
		Likes:    likes,
	}, nil
}

// IsLikeAlreadyDone verify is pretended like is already present
func (post *Post) IsLikeAlreadyDone(like like.Like) bool {
	for _, postLike := range post.Likes {
		if postLike.Equals(like) {
			return true
		}
	}

	return false
}

// LookPresentUserLike find is an user like it has been already given
func (post *Post) LookPresentUserLike(userID string, Type like.Type) like.Like {
	activeState, _ := like.CreateNewLikeState(like.Active)
	for _, postLike := range post.Likes {
		if postLike.UserID == userID && postLike.Type.Equals(Type) && postLike.State.Equals(activeState) {
			return postLike
		}
	}

	return like.Like{}
}
