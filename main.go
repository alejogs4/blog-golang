package main

import (
	"log"
	"net/http"

	"github.com/alejogs4/blog/src/post/application"
	"github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/post/infraestructure/postrepository"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	"github.com/alejogs4/blog/src/shared/infraestructure/token"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
	userrepository "github.com/alejogs4/blog/src/user/infraestructure/userRepository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	err := database.InitDatabase()
	defer database.PostgresDB.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = token.LoadCertificates("certificates/app.rsa", "certificates/app.rsa.pub")
	if err != nil {
		log.Fatal(err)
	}

	// Post use cases
	var postCommands application.PostCommands = application.NewPostCommands(postrepository.NewPostgresRepository(database.PostgresDB))
	var postQueries application.PostQueries = application.NewPostQueries(postrepository.NewPostgresRepository(database.PostgresDB))

	userPostgresReposiry := userrepository.NewUserRepository(database.PostgresDB)

	router := mux.NewRouter()
	router.Handle("/images/{picture}", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	posthttpport.HandlePostHttpRoutes(router, postCommands, postQueries)
	userhttpport.HandleUserRoutes(router, userPostgresReposiry)

	log.Println("Running server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
