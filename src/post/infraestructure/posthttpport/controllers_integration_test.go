package posthttpport_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/post/infraestructure/postrepository"
	"github.com/alejogs4/blog/src/shared/infraestructure/authentication"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	integrationtest "github.com/alejogs4/blog/src/shared/infraestructure/integrationTest"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

var testDatabase *sql.DB

func TestMain(t *testing.M) {
	var err error
	testDatabase, err = database.InitTestDatabase()

	if err != nil {
		log.Fatalf("Error initializing db - %s", err)
		os.Exit(1)
		return
	}

	os.Exit(integrationtest.SetupDatabaseForTesting(t, testDatabase))
}

func TestPostGetAllIntegration(t *testing.T) {
	t.Parallel()

	users, err := integrationtest.PopulateUsers(testDatabase)
	if err != nil {
		t.Errorf("Error: Error inserting users %s", err)
	}

	posts, err := integrationtest.PopulatePosts(users, testDatabase)
	if err != nil {
		t.Errorf("Error: Error inserting posts %s", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	response := httptest.NewRecorder()

	postCommands := application.NewPostCommands(postrepository.NewPostgresRepository(testDatabase))
	postQueries := application.NewPostQueries(postrepository.NewPostgresRepository(testDatabase))

	postsController := posthttpport.NewPostControllers(postCommands, postQueries)
	getAllPostsRouteController := middleware.Chain(postsController.GetAllPostController, httputils.Verb(http.MethodGet))

	getAllPostsRouteController(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Error: Expected status code %d, Received status code %d", http.StatusOK, response.Code)
	}

	var postsInDB struct {
		Posts []post.Post `json:"data"`
	}
	json.NewDecoder(response.Body).Decode(&postsInDB)

	for _, currentPost := range posts {
		postFound := false
		for _, dbPost := range postsInDB.Posts {
			if dbPost.ID == currentPost.ID {
				postFound = true
				break
			}
		}

		if postFound == false {
			t.Errorf("Error: Post %v should have been returned", currentPost)
		}
	}
}

func TestGetByIDIntegration(t *testing.T) {
	t.Parallel()

	users, err := integrationtest.PopulateUsers(testDatabase)
	if err != nil {
		t.Errorf("Error: Error inserting users %s", err)
	}

	posts, err := integrationtest.PopulatePosts(users, testDatabase)
	if err != nil {
		t.Errorf("Error: Error inserting posts %s", err)
	}

	lookPost := posts[0]
	request := httptest.NewRequest(http.MethodGet, "/api/v1/post/{id}", nil)
	response := httptest.NewRecorder()
	withPostIDRequest := mux.SetURLVars(request, map[string]string{"id": lookPost.ID})

	postCommands := application.NewPostCommands(postrepository.NewPostgresRepository(testDatabase))
	postQueries := application.NewPostQueries(postrepository.NewPostgresRepository(testDatabase))

	postsController := posthttpport.NewPostControllers(postCommands, postQueries)
	getAllPostsRouteController := middleware.Chain(postsController.GetPostByIDController, httputils.Verb(http.MethodGet))
	getAllPostsRouteController(response, withPostIDRequest)

	if response.Code != http.StatusOK {
		t.Errorf("Error: Expected status code %d, Received status code %d", http.StatusOK, response.Code)
	}

	var gotResponse struct {
		Post post.Post `json:"data"`
	}
	json.NewDecoder(response.Body).Decode(&gotResponse)

	if gotResponse.Post.ID != lookPost.ID {
		t.Errorf("Error: Expected post %v, Received post %v", lookPost.ID, gotResponse.Post.ID)
	}
}

func TestCreatePostIntegration(t *testing.T) {
	t.Parallel()

	users, _ := integrationtest.PopulateUsers(testDatabase)
	firstUser := users[0]

	loginRequest, loginResponse, loginRoutes := integrationtest.PrepareLoginRequest(firstUser.GetEmail(), firstUser.GetPassword(), testDatabase)
	loginRoutes(loginResponse, loginRequest)

	var loginResponseBody struct {
		Data    userhttpport.LoginResponse `json:"data"`
		Message string                     `json:"message"`
	}

	json.NewDecoder(loginResponse.Body).Decode(&loginResponseBody)
	// Tests petition when is not logged, it's mean that token is not send
	postRepository := postrepository.NewPostgresRepository(testDatabase)
	createResponse, createRequest, postController := preparePostRequest("Title for test purpose", "content", "", postRepository)
	createPostRoute := middleware.Chain(postController.CreatePostController, httputils.Verb(http.MethodPost), authentication.LoginMiddleare())

	createPostRoute(createResponse, createRequest)

	if createResponse.Code != http.StatusUnauthorized {
		t.Errorf("Error: Expected status code %d, Received status code %d", http.StatusUnauthorized, createResponse.Code)
	}

	// Tests petition when is logged, it's mean that token is send
	postTitle := fmt.Sprintf("Title for test purpose %d", rand.Int())
	createLoggedResponse, createLoggedRequest, postLoggedController := preparePostRequest(
		postTitle,
		"content",
		"",
		postRepository,
	)

	createLoggedPostRoute := middleware.Chain(
		postLoggedController.CreatePostController,
		httputils.Verb(http.MethodPost),
		authentication.LoginMiddleare(),
	)

	createLoggedRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResponseBody.Data.Token))
	createLoggedPostRoute(createLoggedResponse, createLoggedRequest)

	if createLoggedResponse.Code != http.StatusCreated {
		t.Errorf("Error: Expected status code %d, Received status code %d", http.StatusCreated, createLoggedResponse.Code)
	}
	// Get all posts
	getAllRequest := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	getAllResponse := httptest.NewRecorder()

	postCommands := application.NewPostCommands(postrepository.NewPostgresRepository(testDatabase))
	postQueries := application.NewPostQueries(postrepository.NewPostgresRepository(testDatabase))

	postsController := posthttpport.NewPostControllers(postCommands, postQueries)
	getAllPostsRouteController := middleware.Chain(postsController.GetAllPostController, httputils.Verb(http.MethodGet))

	getAllPostsRouteController(getAllResponse, getAllRequest)

	var gotResponse struct {
		Posts []post.Post `json:"data"`
	}
	json.NewDecoder(getAllResponse.Body).Decode(&gotResponse)

	foundPost := false
	for _, createdPost := range gotResponse.Posts {
		if createdPost.Title == postTitle {
			foundPost = true
			break
		}
	}

	if foundPost == false {
		t.Errorf("Error: Post with title: %v, it was not created", postTitle)
	}
}
