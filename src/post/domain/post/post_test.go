package post_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

func TestPostEntityUnit(t *testing.T) {
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
		notPresentLike, _ := like.CreateNewLike("third-id", "post-id-3", "user-id-3", like.TLike, like.Active)

		isPresent := post.IsLikeAlreadyDone(searchedLike)
		if !isPresent {
			t.Errorf("Error: Should have found like of user %v and type %v", searchedLike.UserID, searchedLike.Type.GetTypeValue())
		}

		isPresent = post.IsLikeAlreadyDone(notPresentLike)
		if isPresent {
			t.Errorf("Error: Should have not found like %v", notPresentLike)
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

		if err := comment.RemoveComment(); err != nil {
			t.Errorf("Error: error removing comment %s", err)
		}

		if comment.State != post.RemovedComment {
			t.Errorf("Error: comment must have changed to %v", post.RemovedComment)
		}

		err := comment.RemoveComment()
		if err == nil {
			t.Errorf("Error: should have thrown %s", post.ErrInvalidCommentState)
		}
	})
}

func TestPostDTOUnit(t *testing.T) {
	t.Run("Should return a proper instance of PostDTO", func(t *testing.T) {
		rawPost, _ := post.CreateNewPost("123", "123", "title", "content", "picture", []post.Comment{}, []post.Tag{}, []like.Like{})
		postDTO := post.ToPostsDTO(rawPost, 12, 23, 3)
		expectedDTO := post.PostsDTO{
			ID:            "123",
			UserID:        "123",
			Title:         "title",
			Content:       "content",
			Picture:       "picture",
			Likes:         12,
			Dislikes:      23,
			CommentsCount: 3,
		}

		if !reflect.DeepEqual(postDTO, expectedDTO) {
			t.Errorf("Error: expected PostDTO %v, got PostDTO %v", expectedDTO, postDTO)
		}
	})
}

func TestTagEntityUnit(t *testing.T) {
	t.Run("Should throw an error if either id or content are empty", func(t *testing.T) {
		_, err := post.CreateNewTag("", "")
		if err == nil {
			t.Errorf("Error: Should have thrown error %s", post.ErrInvalidTagInformation)
		}

		_, err = post.CreateNewTag("    ", "    ")
		if err == nil {
			t.Errorf("Error: Should have thrown error %s", post.ErrInvalidTagInformation)
		}
	})

	t.Run("Should return proper tag instance if data is correctly provided", func(t *testing.T) {
		tag, _ := post.CreateNewTag("123", "backend")
		expectedTag := post.Tag{
			ID:      "123",
			Content: "backend",
		}

		if !reflect.DeepEqual(tag, expectedTag) {
			t.Errorf("Error: expected tag: %v, got tag: %v", expectedTag, tag)
		}
	})
}
