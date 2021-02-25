package integrationtest

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/shared/infraestructure/token"
	"github.com/alejogs4/blog/src/user/domain/user"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"golang.org/x/crypto/bcrypt"
)

func PopulateUsers(db *sql.DB) ([]user.User, error) {
	newUserOne, _ := user.NewUser(uuid.New().String(), "Alejandro", "Garcia", fake.EmailAddress(), "1234567", true)
	newUserTwo, _ := user.NewUser(uuid.New().String(), "Jose", "Miranda", fake.EmailAddress(), "1234675", true)

	users := []user.User{
		newUserOne,
		newUserTwo,
	}

	for _, currentUser := range users {
		smt, err := db.Prepare("INSERT INTO person(id, firstname, lastname, email, email_verified, password) VALUES($1, $2, $3, $4, $5, $6)")
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

func PopulatePosts(users []user.User, db *sql.DB) ([]post.Post, error) {
	posts := []post.Post{
		{ID: uuid.New().String(), UserID: users[0].GetID(), Title: fake.Title(), Content: "Content 1", Picture: "img/route.jpg", Tags: []post.Tag{}, Comments: []post.Comment{}, Likes: []like.Like{}},
		{ID: uuid.New().String(), UserID: users[1].GetID(), Title: fake.Title(), Content: "Content 2", Picture: "img/route.jpg", Tags: []post.Tag{}, Comments: []post.Comment{}, Likes: []like.Like{}},
	}

	for _, currentPost := range posts {
		smt, err := db.Prepare("INSERT INTO post(id, person_id, title, content, picture) VALUES($1, $2, $3, $4, $5)")
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

func SetupDatabaseForTesting(t *testing.M, db *sql.DB) int {
	enviroment := os.Getenv("ENV")
	if enviroment != "integration_test" {
		return t.Run()
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error: Error closing database %s", err)
			os.Exit(1)
		}
	}()

	if err := token.LoadCertificates("../../../../certificates/app.rsa", "../../../../certificates/app.rsa.pub"); err != nil {
		log.Fatalf("Error: Error initializing token certificates %s", err)
		return 1
	}

	return t.Run()
}
