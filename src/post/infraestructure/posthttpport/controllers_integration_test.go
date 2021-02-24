package posthttpport_test

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/post/infraestructure/postrepository"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	integrationtest "github.com/alejogs4/blog/src/shared/infraestructure/integrationTest"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"

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
