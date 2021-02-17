package user_test

import (
	"fmt"
	"testing"

	"github.com/alejogs4/blog/src/user/domain/user"
)

func TestUserEntity(t *testing.T) {
	t.Run("Should throw an error if any field empty", func(t *testing.T) {
		_, err := user.NewUser(" ", "alejandro", "Garcia", "email-random", "this is my pass", false)
		if err == nil {
			t.Errorf("Error: Should have thrown error %s", user.ErrBadUserData)
		}
	})

	t.Run(fmt.Sprintf("Should return an error if password length is less than %d", user.MinPasswordLength), func(t *testing.T) {
		_, err := user.NewUser("id", "alejandro", "Garcia", "email-random", "123", false)
		if err == nil {
			t.Errorf("Error: Should have thrown error %s", user.ErrTooShortUserPassword)
		}
	})
}
