package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

var (
	host = os.Getenv("DATABASE_HOST")
	port = 5432
)

var (
	// PostgresDB database connection to interact with the database
	PostgresDB *sql.DB
	once       sync.Once
)

func InitDatabase() error {
	var err error
	once.Do(func() {
		err = initDatabase()
	})

	return err
}

func InitTestDatabase() (*sql.DB, error) {
	enviroment := os.Getenv("ENV")
	if enviroment != "integration_test" {
		return nil, nil
	}

	databaseUser := os.Getenv("POSTGRES_USER")
	databaseName := os.Getenv("POSTGRES_DB")

	postgresInformation := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, databaseUser, databaseName)

	testDB, err := sql.Open("postgres", postgresInformation)

	if err != nil {
		return nil, err
	}

	err = testDB.Ping()
	if err != nil {
		return nil, err
	}

	return testDB, nil
}

func initDatabase() error {
	databaseUser := os.Getenv("POSTGRES_USER")
	databaseName := os.Getenv("POSTGRES_DB")

	postgresInformation := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, databaseUser, databaseName)

	var err error
	PostgresDB, err = sql.Open("postgres", postgresInformation)

	if err != nil {
		return err
	}

	err = PostgresDB.Ping()
	if err != nil {
		return err
	}

	return err
}
