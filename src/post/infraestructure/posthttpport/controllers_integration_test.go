package posthttpport_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	posthttppost "github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/post/infraestructure/postrepository"
	"github.com/alejogs4/blog/src/shared/infraestructure/authentication"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	integrationtest "github.com/alejogs4/blog/src/shared/infraestructure/integrationTest"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
	"github.com/gorilla/mux"
	"github.com/icrowley/fake"

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

	response, request, getAllPostsRouteController := prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllPostsRouteController(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Error: Expected status code %d, Received status code %d", http.StatusOK, response.Code)
	}

	var postsInDB struct {
		Posts []post.PostsDTO `json:"data"`
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
	postTitle := fake.Title()
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
	getAllResponse, getAllRequest, getAllPostsRouteController := prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllPostsRouteController(getAllResponse, getAllRequest)

	var gotResponse struct {
		Posts []post.PostsDTO `json:"data"`
	}
	json.NewDecoder(getAllResponse.Body).Decode(&gotResponse)

	foundPost := existPost(func(createdPost post.PostsDTO) bool {
		return createdPost.Title == postTitle
	}, gotResponse.Posts)

	if foundPost == false {
		t.Errorf("Error: Post with title: %v, it was not created", postTitle)
	}
}

func TestAddLikeIntegration(t *testing.T) {
	t.Parallel()

	users, err := integrationtest.PopulateUsers(testDatabase)
	if err != nil {
		t.Errorf("Error: error inserting users %s", err)
	}

	posts, err := integrationtest.PopulatePosts(users, testDatabase)
	if err != nil {
		t.Errorf("Error: error inserting posts %s", err)
	}

	firstUser := users[0]

	loginRequest, loginResponse, loginHandler := integrationtest.PrepareLoginRequest(
		firstUser.GetEmail(),
		firstUser.GetPassword(),
		testDatabase,
	)

	loginHandler(loginResponse, loginRequest)
	var loginResponseBody struct {
		Data    userhttpport.LoginResponse `json:"data"`
		Message string                     `json:"message"`
	}

	err = json.NewDecoder(loginResponse.Body).Decode(&loginResponseBody)
	if err != nil {
		t.Errorf("Error: error deconding response body %s", err)
	}

	firstPost := posts[0]
	postgresRepository := postrepository.NewPostgresRepository(testDatabase)
	addLikeResponse, unloggedAddLikeRequest, unLoggedAddLikeHandler := prepareAddLikeRequest(like.TLike, firstPost.ID, postgresRepository)

	addLikeHandler := middleware.Chain(
		unLoggedAddLikeHandler.AddPostLikeController,
		httputils.Verb(http.MethodPost),
		authentication.LoginMiddleare(),
	)
	addLikeHandler(addLikeResponse, unloggedAddLikeRequest)

	if addLikeResponse.Code != http.StatusUnauthorized {
		t.Errorf("Error: Expected status code %d, Received status code %d", http.StatusUnauthorized, addLikeResponse.Code)
	}

	addLikeResponse, loggedAddLikeRequest, _ := prepareAddLikeRequest(like.TLike, firstPost.ID, postgresRepository)
	loggedAddLikeRequest.Header.Set("Authorization", "Bearer "+loginResponseBody.Data.Token)
	addLikeHandler(addLikeResponse, loggedAddLikeRequest)

	if addLikeResponse.Code != http.StatusCreated {
		t.Errorf("Error: Expected status code %d, Received status code %d", http.StatusCreated, addLikeResponse.Code)
	}

	// Get post
	getAllResponse, getAllRequest, getAllHandler := prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllHandler(getAllResponse, getAllRequest)

	var gotResponse struct {
		Posts []post.PostsDTO `json:"data"`
	}
	json.NewDecoder(getAllResponse.Body).Decode(&gotResponse)

	foundLikedPost := existPost(func(storedPost post.PostsDTO) bool {
		return storedPost.ID == firstPost.ID && storedPost.Likes == 1
	}, gotResponse.Posts)

	if foundLikedPost == false {
		t.Errorf("Error: Post with id %v and liked was not found", firstPost.ID)
	}

	// New like
	addLikeResponse, loggedRemoveLikeRequest, _ := prepareAddLikeRequest(like.TLike, firstPost.ID, postgresRepository)
	loggedRemoveLikeRequest.Header.Set("Authorization", "Bearer "+loginResponseBody.Data.Token)
	addLikeHandler(addLikeResponse, loggedRemoveLikeRequest)

	getAllResponse, getAllRequest, getAllHandler = prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllHandler(getAllResponse, getAllRequest)

	json.NewDecoder(getAllResponse.Body).Decode(&gotResponse)

	foundRemovedPostLike := existPost(func(storedPost post.PostsDTO) bool {
		return storedPost.ID == firstPost.ID && storedPost.Likes == 0
	}, gotResponse.Posts)

	if foundRemovedPostLike == false {
		t.Errorf("Error: Post with id %v and removed like was not found", firstPost.ID)
	}

	// Dislike
	addLikeResponse, loggedDislikeRequest, _ := prepareAddLikeRequest(like.Dislike, firstPost.ID, postgresRepository)
	loggedDislikeRequest.Header.Set("Authorization", "Bearer "+loginResponseBody.Data.Token)
	addLikeHandler(addLikeResponse, loggedDislikeRequest)

	getAllResponse, getAllRequest, getAllHandler = prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllHandler(getAllResponse, getAllRequest)

	json.NewDecoder(getAllResponse.Body).Decode(&gotResponse)

	foundDislikePost := existPost(func(storedPost post.PostsDTO) bool {
		return storedPost.ID == firstPost.ID && storedPost.Likes == 0 && storedPost.Dislikes == 1
	}, gotResponse.Posts)

	if foundDislikePost == false {
		t.Errorf("Error: Post with id %v and disliked was not found", firstPost.ID)
	}

	// Like again, this should remove the last dislike and add a like
	addLikeResponse, loggedLikeRequest, _ := prepareAddLikeRequest(like.TLike, firstPost.ID, postgresRepository)
	loggedLikeRequest.Header.Set("Authorization", "Bearer "+loginResponseBody.Data.Token)
	addLikeHandler(addLikeResponse, loggedLikeRequest)

	getAllResponse, getAllRequest, getAllHandler = prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllHandler(getAllResponse, getAllRequest)
	json.NewDecoder(getAllResponse.Body).Decode(&gotResponse)

	foundLikedPost = existPost(func(storedPost post.PostsDTO) bool {
		return storedPost.ID == firstPost.ID && storedPost.Likes == 1 && storedPost.Dislikes == 0
	}, gotResponse.Posts)

	if foundLikedPost == false {
		t.Errorf("Error: Post with id %v and liked again was not found", firstPost.ID)
	}
}

func TestAddRemoveCommentIntegration(t *testing.T) {
	t.Parallel()

	users, err := integrationtest.PopulateUsers(testDatabase)
	if err != nil {
		t.Errorf("Error: error populating users %s", err)
	}

	posts, err := integrationtest.PopulatePosts(users, testDatabase)
	if err != nil {
		t.Errorf("Error: error populating posts %s", err)
	}

	loggedUser := users[0]
	loginRequest, loginResponse, loginController := integrationtest.PrepareLoginRequest(
		loggedUser.GetEmail(),
		loggedUser.GetPassword(),
		testDatabase,
	)

	loginController(loginResponse, loginRequest)
	if loginResponse.Code != http.StatusOK {
		t.Errorf("Error: expected status code %d, received status code %d", http.StatusOK, loginResponse.Code)
	}

	var loginInformation struct {
		Data userhttpport.LoginResponse `json:"data"`
	}
	json.NewDecoder(loginResponse.Body).Decode(&loginInformation)

	usedPost := posts[0]
	commentContent := []byte(fmt.Sprintf(`{"content": "%v"}`, fake.ParagraphsN(1)))

	addCommentRequest := httptest.NewRequest(http.MethodPost, "/api/v1/post/{id}/comment", bytes.NewBuffer(commentContent))
	addCommentResponse := httptest.NewRecorder()

	addCommentRequest.Header.Set("Authorization", "Bearer "+loginInformation.Data.Token)
	withPostIDRequest := mux.SetURLVars(addCommentRequest, map[string]string{"id": usedPost.ID})

	postgresRepository := postrepository.NewPostgresRepository(testDatabase)
	postCommands := application.NewPostCommands(postgresRepository)
	postQueries := application.NewPostQueries(postgresRepository)

	addPostCommentController := posthttppost.NewPostControllers(postCommands, postQueries).AddPostComment

	addPostCommentRoute := middleware.Chain(
		addPostCommentController,
		httputils.Verb(http.MethodPost),
		authentication.LoginMiddleare(),
	)

	addPostCommentRoute(addCommentResponse, withPostIDRequest)
	if addCommentResponse.Code != http.StatusCreated {
		t.Errorf("Error: expected status code %d, received status code %d", http.StatusCreated, addCommentResponse.Code)
	}

	var commentID struct {
		Data struct {
			CommentID string `json:"comment_id"`
		} `json:"data"`
	}
	json.NewDecoder(addCommentResponse.Body).Decode(&commentID)

	var gotPostsResponse struct {
		Posts []post.PostsDTO `json:"data"`
	}
	getAllResponse, getAllRequest, getAllHandler := prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllHandler(getAllResponse, getAllRequest)
	json.NewDecoder(getAllResponse.Body).Decode(&gotPostsResponse)

	oneCommentPost := existPost(func(storedPost post.PostsDTO) bool {
		return storedPost.CommentsCount == 1 && storedPost.ID == usedPost.ID
	}, gotPostsResponse.Posts)

	if !oneCommentPost {
		t.Errorf("Error: Post with ID %v should have one comment", usedPost.ID)
	}

	removeCommentRequest := httptest.NewRequest(http.MethodDelete, "/api/v1/comment/{id}", nil)
	removeCommentResponse := httptest.NewRecorder()
	removeCommentRequest.Header.Set("Authorization", "Bearer "+loginInformation.Data.Token)
	withPostIDRequest = mux.SetURLVars(removeCommentRequest, map[string]string{"id": commentID.Data.CommentID})

	removeCommentController := posthttppost.NewPostControllers(postCommands, postQueries).RemoveComment
	removeCommentRoute := middleware.Chain(
		removeCommentController,
		httputils.Verb(http.MethodDelete),
		authentication.LoginMiddleare(),
	)

	removeCommentRoute(removeCommentResponse, withPostIDRequest)
	if removeCommentResponse.Code != http.StatusOK {
		t.Errorf("Error: expected status code %d, received status code %d", http.StatusOK, removeCommentResponse.Code)
	}

	getAllResponse, getAllRequest, getAllHandler = prepareGetAllPostsRequest(postrepository.NewPostgresRepository(testDatabase))
	getAllHandler(getAllResponse, getAllRequest)
	json.NewDecoder(getAllResponse.Body).Decode(&gotPostsResponse)

	noCommentPost := existPost(func(storedPost post.PostsDTO) bool {
		return storedPost.CommentsCount == 0 && storedPost.ID == usedPost.ID
	}, gotPostsResponse.Posts)

	if !noCommentPost {
		t.Errorf("Error: Post with ID %v should have no comment", usedPost.ID)
	}
}
