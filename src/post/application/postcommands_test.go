package application_test

import (
	"errors"
	"testing"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

type PostRepositoryMock struct {
	RemovedLike  *like.Like
	AddedLike    *like.Like
	ReturnedPost post.Post
}

func (pr PostRepositoryMock) CreatePost(post post.Post) error {
	return nil
}

func (pr PostRepositoryMock) AddLike(postID string, like like.Like) error {
	*pr.AddedLike = like
	return nil
}

func (pr PostRepositoryMock) RemoveLike(like like.Like) error {
	*pr.RemovedLike = like
	return nil
}

func (pr PostRepositoryMock) AddComment(comment post.Comment) error {
	return nil
}

func (pr PostRepositoryMock) RemoveComment(comment post.Comment) error {
	return nil
}

func (pr PostRepositoryMock) GetPostCommentByID(id string) (post.Comment, error) {
	return post.Comment{}, nil
}

func (pr PostRepositoryMock) GetAllPosts() ([]post.PostsDTO, error) {
	return []post.PostsDTO{}, nil
}

func (pr PostRepositoryMock) GetPostLikes(postID string) ([]like.Like, error) {
	return []like.Like{}, nil
}

func (pr PostRepositoryMock) GetPostByID(postID string) (post.Post, error) {
	return pr.ReturnedPost, nil
}

func TestPostCommandsCreateNewPost(t *testing.T) {
	postCommands := application.NewPostCommands(PostRepositoryMock{})

	t.Run("Should return nil if all fields are correctly provided", func(t *testing.T) {
		err := postCommands.CreateNewPost("123", "title", "content", "picture", []post.Tag{})

		if err != nil {
			t.Errorf("Error: Nil error was expected, received %v", err)
		}
	})

	t.Run("Should return an error if fields are empty", func(t *testing.T) {
		err := postCommands.CreateNewPost("123", "   ", "content", "picture", []post.Tag{})

		if err == nil {
			t.Errorf("Error: Expected error %v, received nil", post.ErrBadPostContent)
		}

		if !errors.Is(err, post.ErrBadPostContent) {
			t.Errorf("Error: Expected error %v, received %v", post.ErrBadPostContent, err)
		}
	})
}

func TestPostCommandsCreateNewComment(t *testing.T) {
	postCommands := application.NewPostCommands(PostRepositoryMock{})
	t.Run("Should return ErrBadCommentContent if not all fields were provided", func(t *testing.T) {
		err := postCommands.CreateNewComment("  ", " ", "content")

		if err == nil {
			t.Errorf("Error: expected error %v, received nil", post.ErrBadCommentContent)
		}

		if !errors.Is(err, post.ErrBadCommentContent) {
			t.Errorf("Error: expected error %v, received error %v", post.ErrBadCommentContent, err)
		}
	})

	t.Run("Should return nil if all fields are provided", func(t *testing.T) {
		err := postCommands.CreateNewComment("user-id", "post-id", "content")

		if err != nil {
			t.Errorf("Error: expected error nil, received  %v", err)
		}
	})
}

func TestPostCommandsAddLike(t *testing.T) {
	t.Run("Should return an error if like type is not either like or dislike", func(t *testing.T) {
		postCommands := application.NewPostCommands(PostRepositoryMock{})
		err := postCommands.AddLike("user-id", "post-id", "invalid-like-type")

		if err == nil {
			t.Errorf("Error: expected error was %v, received nil instead", like.ErrInvalidLikeType)
		}

		if !errors.Is(err, like.ErrInvalidLikeType) {
			t.Errorf("Error: expected error was %v, received %v", like.ErrInvalidLikeType, err)
		}
	})

	t.Run("Should remove like if the same like type is present", func(t *testing.T) {
		targetLike := like.Like{ID: "123", PostID: "123", UserID: "user-id", Type: like.Type{Value: like.TLike}, State: like.State{Value: like.Active}}
		removedLike := like.Like{}
		mockPostRepository := PostRepositoryMock{
			ReturnedPost: post.Post{
				ID:       "123",
				UserID:   "user-id",
				Title:    "title",
				Content:  "content",
				Picture:  "picture",
				Tags:     []post.Tag{},
				Comments: []post.Comment{},
				Likes: []like.Like{
					targetLike,
				},
			},
			RemovedLike: &removedLike,
		}
		postCommands := application.NewPostCommands(mockPostRepository)

		err := postCommands.AddLike(targetLike.UserID, targetLike.PostID, targetLike.Type.Value)
		if err != nil {
			t.Errorf("Error: expected error nil was expected, received %v", err)
		}

		if !removedLike.Equals(targetLike) {
			t.Errorf("Error: RemoveLike was expected to had removed %v", targetLike)
		}
	})

	t.Run("Should remove previously given but opposite like if a new one is added and add this new one", func(t *testing.T) {
		oppositeLike := like.Like{ID: "122", PostID: "123", UserID: "user-id", Type: like.Type{Value: like.Dislike}, State: like.State{Value: like.Active}}
		targetLike := like.Like{ID: "123", PostID: "123", UserID: "user-id", Type: like.Type{Value: like.TLike}, State: like.State{Value: like.Active}}

		removedLike := like.Like{}
		addedLike := like.Like{}

		mockPostRepository := PostRepositoryMock{
			ReturnedPost: post.Post{
				ID:       "123",
				UserID:   "user-id",
				Title:    "title",
				Content:  "content",
				Picture:  "picture",
				Tags:     []post.Tag{},
				Comments: []post.Comment{},
				Likes: []like.Like{
					oppositeLike,
				},
			},
			RemovedLike: &removedLike,
			AddedLike:   &addedLike,
		}
		postCommands := application.NewPostCommands(mockPostRepository)
		err := postCommands.AddLike(targetLike.UserID, targetLike.PostID, targetLike.Type.Value)

		if err != nil {
			t.Errorf("Error: expected error nil was expected, received %v", err)
		}

		if !removedLike.Equals(oppositeLike) {
			t.Errorf("Error: RemoveLike was expected to had removed %v", oppositeLike)
		}

		if !addedLike.Equals(targetLike) {
			t.Errorf("Error: AddLike was expected to had added like %v", targetLike)
		}
	})
}
