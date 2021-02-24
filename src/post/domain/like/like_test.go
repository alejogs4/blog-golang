package like_test

import (
	"errors"
	"testing"

	"github.com/alejogs4/blog/src/post/domain/like"
)

func TestLikeEntityUnit(t *testing.T) {
	t.Run("Should throw an error if right content is not provided", func(t *testing.T) {
		_, err := like.CreateNewLike("", "post-id", "user-id", like.Dislike, like.Active)
		if err == nil {
			t.Errorf("Error: Should have thrown error -- %s", like.ErrBadLikeContent)
		}

		_, err = like.CreateNewLike("", "post-id", "user-id", "invalid-type", like.Active)
		if err == nil {
			t.Errorf("Error: Should have thrown error -- %s", like.ErrInvalidLikeType)
		}

		_, err = like.CreateNewLike("id", "post-id", "user-id", like.Dislike, "invalid-state")
		if err == nil || !errors.Is(err, like.ErrInvalidLikeState) {
			t.Errorf("Error: Should have thrown error -- %s", like.ErrInvalidLikeState)
		}
	})

	t.Run("Should switch like type to the opposite of the current", func(t *testing.T) {
		currentLike, _ := like.CreateNewLike("id", "post-id", "user-id", like.Dislike, like.Active)

		if currentLike.Type.GetTypeValue() != like.Dislike {
			t.Errorf("Error: At the beginning should be type: %v got: %v", like.Dislike, currentLike.Type.GetTypeValue())
		}

		newLike, _ := currentLike.SwitchType()
		if newLike.Type.GetTypeValue() != like.TLike {
			t.Errorf("Error: Switched like should be: %v got: %v", like.TLike, newLike.Type.GetTypeValue())
		}

		dislike, _ := newLike.SwitchType()
		if dislike.Type.GetTypeValue() != like.Dislike {
			t.Errorf("Error: Switched like should be: %v got: %v", like.Dislike, dislike.Type.GetTypeValue())
		}
	})

	t.Run("Should return true if both like are equal", func(t *testing.T) {
		firstLike, _ := like.CreateNewLike("id", "postid", "userid", like.TLike, like.Active)
		secondLike, _ := like.CreateNewLike("id", "postid", "userid", like.TLike, like.Active)

		got := firstLike.Equals(secondLike)
		want := true
		if got != want {
			t.Errorf("Error: Should have returned %v but got %v", want, got)
		}
	})

	t.Run("Should return false if both like are not equal", func(t *testing.T) {
		firstLike, _ := like.CreateNewLike("id", "postid", "user-id-1", like.TLike, like.Active)
		secondLike, _ := like.CreateNewLike("id", "postid-2", "user-id-2", like.Dislike, like.Active)

		got := firstLike.Equals(secondLike)
		want := false
		if got != want {
			t.Errorf("Error: Should have returned %v but got %v", want, got)
		}
	})
}
