package main

import (
	"log"
	"net/http"

	"github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
	_ "github.com/lib/pq"
)

func main() {
	err := database.InitDatabase()
	defer database.PostgresDB.Close()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	posthttpport.HandlePostHttpRoutes(router)
	userhttpport.HandleUserRoutes(router)

	log.Println("Running server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
