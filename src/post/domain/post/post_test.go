package post_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

func TestPostEntity(t *testing.T) {
	t.Run("Post create should throws error if is created with empty fields", func(t *testing.T) {
		_, err := post.CreateNewPost("", "", "any title", "content", "picture", []post.Comment{}, []post.Tag{}, []like.Like{})
		if err == nil {
			t.Errorf("Error: It should have thrown error -- %s", post.ErrBadPostContent)
		}

		_, err = post.CreateNewPost("id", "", "any title", "content", "picture", []post.Comment{}, []post.Tag{}, []like.Like{})
		if err == nil {
			t.Errorf("Error: It should have thrown error -- %s", post.ErrBadPostContent)
		}

		_, err = post.CreateNewPost("id", "userid", "", "content", "picture", []post.Comment{}, []post.Tag{}, []like.Like{})
		if err == nil {
			t.Errorf("Error: It should have thrown error -- %s", post.ErrBadPostContent)
		}

		_, err = post.CreateNewPost("id", "userid", "any title", "", "picture", []post.Comment{}, []post.Tag{}, []like.Like{})
		if err == nil {
			t.Errorf("Error: It should have thrown error -- %s", post.ErrBadPostContent)
		}

		_, err = post.CreateNewPost("id", "userid", "any title", "picture", "", []post.Comment{}, []post.Tag{}, []like.Like{})
		if err == nil {
			t.Errorf("Error: It should have thrown error -- %s", post.ErrBadPostContent)
		}
	})

	t.Run("Should return true if like it has been already done, false otherwise", func(t *testing.T) {
		firstLike, _ := like.CreateNewLike("first-id", "post-id", "user-id", like.TLike, like.Active)
		secondLike, _ := like.CreateNewLike("second-id", "post-id", "user-id", like.Dislike, like.Active)

		post, _ := post.CreateNewPost("post-id", "user-id", "any title", "content", "picture", []post.Comment{}, []post.Tag{}, []like.Like{
			firstLike,
			secondLike,
		})

		searchedLike, _ := like.CreateNewLike("second-id", "post-id", "user-id", like.Dislike, like.Active)
		isPresent := post.IsLikeAlreadyDone(searchedLike)
		if !isPresent {
			t.Errorf("Error: Should have found like of user %v and type %v", searchedLike.UserID, searchedLike.Type.GetTypeValue())
		}
	})

	t.Run("Should return present user like given like type and user id", func(t *testing.T) {
		secondLikeUserID := "user-id-2"
		secondLikeType, _ := like.CreateNewLikeType(like.Dislike)

		firstLike, _ := like.CreateNewLike("first-id", "post-id", "user-id-1", like.TLike, like.Active)
		secondLike, _ := like.CreateNewLike("second-id", "post-id", secondLikeUserID, secondLikeType.GetTypeValue(), like.Active)

		post, _ := post.CreateNewPost("post-id", "user-id", "any title", "content", "picture", []post.Comment{}, []post.Tag{}, []like.Like{
			firstLike,
			secondLike,
		})

		foundLike := post.LookPresentUserLike(secondLikeUserID, secondLikeType)
		if foundLike.UserID != secondLikeUserID || !foundLike.Type.Equals(secondLikeType) {
			t.Errorf("Error: Should have found a like with user id %v and type %v", foundLike.UserID, foundLike.Type.GetTypeValue())
		}
	})
}

func TestCommentEntity(t *testing.T) {
	t.Run("Should throw an error if a field is empty", func(t *testing.T) {
		_, err := post.CreateNewComment("id", "     ", "      ", "content")
		if err == nil {
			t.Errorf("Error: Should have thrown the error -- %s", post.ErrBadCommentContent)
		}
	})

	t.Run("Should not throw an error if all fields are filled", func(t *testing.T) {
		_, err := post.CreateNewComment("id", "post-id", "user-id", "content")
		if err != nil {
			t.Errorf("Error: Should not have thrown the error -- %s", err.Error())
		}
	})

	t.Run(fmt.Sprintf("Should throw an error if comment length is greather than %d characters", post.MaxCommentContent), func(t *testing.T) {
		var longContent bytes.Buffer
		for i := 0; i < post.MaxCommentContent+10; i++ {
			longContent.WriteString("a")
		}

		_, err := post.CreateNewComment("id", "dd", "dd", longContent.String())
		if err == nil {
			t.Errorf("Error: Should have thrown the error -- %s", post.ErrInvalidCommentLength)
		}

		var shortContent bytes.Buffer
		for i := 0; i < post.MaxCommentContent-10; i++ {
			shortContent.WriteString("a")
		}

		_, err = post.CreateNewComment("id", "dd", "dd", shortContent.String())
		if err != nil {
			t.Errorf("Error: Should not have thrown the error -- %s", err.Error())
		}
	})

	t.Run("Should change comment state to removed if RemoveComment method is executed", func(t *testing.T) {
		comment, _ := post.CreateNewComment("id", "post-id", "user-id", "content")
		if comment.State != post.ActiveComment {
			t.Errorf("Error: comment state should be initially %v", post.ActiveComment)
		}

		comment.RemoveComment()
		if comment.State != post.RemovedComment {
			t.Errorf("Error: comment must have changed to %v", post.RemovedComment)
		}

		err := comment.RemoveComment()
		if err == nil {
			t.Errorf("Error: should have thrown %s", post.ErrInvalidCommentState)
		}
	})
}
