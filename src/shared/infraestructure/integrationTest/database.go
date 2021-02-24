package integrationtest

import (
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	"github.com/alejogs4/blog/src/user/domain/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func PopulateUsers() ([]user.User, error) {
	newUserOne, _ := user.NewUser(uuid.New().String(), "Alejandro", "Garcia", "alejogs4@gmail.com", "12345", true)
	newUserTwo, _ := user.NewUser(uuid.New().String(), "Jose", "Miranda", "josemiranda@gmail.com", "12345", true)
	newUserThree, _ := user.NewUser(uuid.New().String(), "Miguel", "Velez", "miguelito99@gmail.com", "12345", true)
	newUserFour, _ := user.NewUser(uuid.New().String(), "Mauricio", "Brunal Mestra", "mgbrunal@gmail.com", "12345", true)

	users := []user.User{
		newUserOne,
		newUserTwo,
		newUserThree,
		newUserFour,
	}

	for _, currentUser := range users {
		smt, err := database.PostgresDB.Prepare("INSERT INTO person(id, firstname, lastname, email, email_verified, password) VALUES($1, $2, $3, $4, $5, $6)")
		if err != nil {
			return []user.User{}, err
		}

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(currentUser.GetPassword()), 14)
		if err != nil {
			return []user.User{}, err
		}

		_, err = smt.Exec(
			currentUser.GetID(),
			currentUser.GetFirstname(),
			currentUser.GetLastname(),
			currentUser.GetEmail(),
			currentUser.GetEmailVerified(),
			encryptedPassword,
		)
		if err != nil {
			return []user.User{}, err
		}
	}

	return users, nil
}

func PopulatePosts(users []user.User) ([]post.Post, error) {
	posts := []post.Post{
		{ID: uuid.New().String(), UserID: users[0].GetID(), Title: "Title 1", Content: "Content 1", Picture: "img/route.jpg", Tags: []post.Tag{}, Comments: []post.Comment{}, Likes: []like.Like{}},
		{ID: uuid.New().String(), UserID: users[1].GetID(), Title: "Title 2", Content: "Content 2", Picture: "img/route.jpg", Tags: []post.Tag{}, Comments: []post.Comment{}, Likes: []like.Like{}},
		{ID: uuid.New().String(), UserID: users[1].GetID(), Title: "Title 3", Content: "Content 3", Picture: "img/route.jpg", Tags: []post.Tag{}, Comments: []post.Comment{}, Likes: []like.Like{}},
	}

	for _, currentPost := range posts {
		smt, err := database.PostgresDB.Prepare("INSERT INTO post(id, person_id, title, content, picture) VALUES($1, $2, $3, $4, $5)")
		if err != nil {
			return []post.Post{}, err
		}

		_, err = smt.Exec(currentPost.ID, currentPost.UserID, currentPost.Title, currentPost.Content, currentPost.Picture)
		if err != nil {
			return []post.Post{}, err
		}
	}

	return posts, nil
}

func TruncateDatabase() error {
	_, err := database.PostgresDB.Exec("TRUNCATE TABLE person, post, comment, post_like, tag, post_tag")
	if err != nil {
		return err
	}

	return nil
}
