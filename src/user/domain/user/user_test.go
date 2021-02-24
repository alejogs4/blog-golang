package user_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alejogs4/blog/src/user/domain/user"
)

func TestUserEntityUnit(t *testing.T) {
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

	t.Run("Should return expected user if all validations were passed", func(t *testing.T) {
		expectedUserID := "123"
		gotUser, _ := user.NewUser(expectedUserID, "alejandro", "Garcia", "email-random", "this is my pass", false)

		if expectedUserID != gotUser.GetID() {
			t.Errorf("Error: Expected user id: %v, Got user id: %v", expectedUserID, gotUser.GetID())
		}
	})

	t.Run("Should create proper user dto", func(t *testing.T) {
		createdUser, _ := user.NewUser("123", "Alejandro", "Garcia", "alejogs4@gmail.com", "123345667", true)
		userDTO := user.ToDTO(createdUser)
		expectedDTO := user.UserDTO{
			ID:            "123",
			Firstname:     "Alejandro",
			Lastname:      "Garcia",
			Email:         "alejogs4@gmail.com",
			EmailVerified: true,
		}

		if !reflect.DeepEqual(userDTO, expectedDTO) {
			t.Errorf("Error: expected %v, got %v", expectedDTO, userDTO)
		}
	})
}
