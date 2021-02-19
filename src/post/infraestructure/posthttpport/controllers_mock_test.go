package posthttpport_test

import (
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

type mockPostRepositoryOK struct{}

func (mock mockPostRepositoryOK) CreatePost(newPost post.Post) error {
	return nil
}

func (mock mockPostRepositoryOK) AddLike(postID string, like like.Like) error {
	return nil
}

func (mock mockPostRepositoryOK) AddComment(comment post.Comment) error {
	return nil
}

func (mock mockPostRepositoryOK) RemoveComment(comment post.Comment) error {
	return nil
}

func (mock mockPostRepositoryOK) RemoveLike(like like.Like) error {
	return nil
}

func (mock mockPostRepositoryOK) GetAllPosts() ([]post.PostsDTO, error) {
	return []post.PostsDTO{}, nil
}

func (mock mockPostRepositoryOK) GetPostLikes(postID string) ([]like.Like, error) {
	return []like.Like{}, nil
}

func (mock mockPostRepositoryOK) GetPostByID(postID string) (post.Post, error) {
	return post.Post{}, nil
}

func (mock mockPostRepositoryOK) GetPostCommentByID(id string) (post.Comment, error) {
	return post.Comment{}, nil
}
