package application_test

import (
	"errors"
	"testing"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

func TestPostCommandsCreateNewPostUnit(t *testing.T) {
	postCommands := application.NewPostCommands(postRepositoryMock{})

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

func TestPostCommandsCreateNewCommentUnit(t *testing.T) {
	postCommands := application.NewPostCommands(postRepositoryMock{})
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

func TestPostCommandsAddLikeUnit(t *testing.T) {
	t.Run("Should return an error if like type is not either like or dislike", func(t *testing.T) {
		postCommands := application.NewPostCommands(postRepositoryMock{})
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
		mockPostRepository := postRepositoryMock{
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

		mockPostRepository := postRepositoryMock{
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

func TestRemovePostCommentUnit(t *testing.T) {
	t.Run("Should return error if user to remove the post is no the post owner", func(t *testing.T) {
		mockPostRepository := postRepositoryMock{
			ReturnedComment: post.Comment{
				UserID: "1234567",
			},
		}

		postCommands := application.NewPostCommands(mockPostRepository)
		err := postCommands.RemovePostComment("id", "different-user-id")

		if err == nil {
			t.Errorf("Error: Expected error %v, received nil", post.ErrNoCommentOwner)
		}

		if !errors.Is(err, post.ErrNoCommentOwner) {
			t.Errorf("Error: Expected error %v, received %v", post.ErrNoCommentOwner, err)
		}
	})

	t.Run("Should return no error if user who remove the post is the post owner", func(t *testing.T) {
		userID := "1234567"
		mockPostRepository := postRepositoryMock{
			ReturnedComment: post.Comment{
				UserID: userID,
			},
		}

		postCommands := application.NewPostCommands(mockPostRepository)
		err := postCommands.RemovePostComment("id", userID)

		if err != nil {
			t.Errorf("Error: No error was expected, received %v", err)
		}
	})
}
