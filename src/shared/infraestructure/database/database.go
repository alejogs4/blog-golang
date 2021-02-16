package database

import (
	"database/sql"
	"fmt"
	"sync"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "blog"
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

func initDatabase() error {
	postgresInformation := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, user, dbname)

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
