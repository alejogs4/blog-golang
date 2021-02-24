package integrationtest_test

import (
	"log"
	"os"
	"testing"

	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	integrationtest "github.com/alejogs4/blog/src/shared/infraestructure/integrationTest"
	_ "github.com/lib/pq"
)

func testMain(t *testing.M) int {
	enviroment := os.Getenv("ENV")

	defer func() {
		if enviroment == "integration_test" {
			err := integrationtest.TruncateDatabase()
			if err != nil {
				log.Fatalf("Error: Error truncating database %s", err)
				os.Exit(1)
			}

			if err := database.PostgresDB.Close(); err != nil {
				log.Fatalf("Error: Error closing database %s", err)
				os.Exit(1)
			}
		}
	}()

	if enviroment == "integration_test" {
		if err := database.InitDatabase(); err != nil {
			log.Fatalf("Error: Error initializing database %s", err)
			return 1
		}
	}

	return t.Run()
}

func TestMain(t *testing.M) {
	os.Exit(testMain(t))
}
