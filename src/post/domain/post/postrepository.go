package post

import "github.com/alejogs4/blog/src/post/domain/like"

type PostRepository interface {
	CreatePost(post Post) error
	AddLike(postID string, like like.Like) error
	RemoveLike(like like.Like) error
	AddComment(comment Comment) error
	RemoveComment(comment Comment) error

	GetPostCommentByID(id string) (Comment, error)
	GetAllPosts() ([]PostsDTO, error)
	GetPostLikes(postID string) ([]like.Like, error)
	GetPostByID(postID string) (Post, error)
}
