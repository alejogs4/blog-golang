package application_test

import (
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

type postRepositoryMock struct {
	RemovedLike     *like.Like
	AddedLike       *like.Like
	ReturnedPost    post.Post
	ReturnedComment post.Comment
}

func (pr postRepositoryMock) CreatePost(post post.Post) error {
	return nil
}

func (pr postRepositoryMock) AddLike(postID string, like like.Like) error {
	*pr.AddedLike = like
	return nil
}

func (pr postRepositoryMock) RemoveLike(like like.Like) error {
	*pr.RemovedLike = like
	return nil
}

func (pr postRepositoryMock) AddComment(comment post.Comment) error {
	return nil
}

func (pr postRepositoryMock) RemoveComment(comment post.Comment) error {
	return nil
}

func (pr postRepositoryMock) GetPostCommentByID(id string) (post.Comment, error) {
	return pr.ReturnedComment, nil
}

func (pr postRepositoryMock) GetAllPosts() ([]post.PostsDTO, error) {
	return []post.PostsDTO{}, nil
}

func (pr postRepositoryMock) GetPostLikes(postID string) ([]like.Like, error) {
	return []like.Like{}, nil
}

func (pr postRepositoryMock) GetPostByID(postID string) (post.Post, error) {
	return pr.ReturnedPost, nil
}
