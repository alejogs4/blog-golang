package main

import (
	"log"
	"net/http"

	"github.com/alejogs4/blog/src/post/infraestructure/posthttpport"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	"github.com/alejogs4/blog/src/shared/infraestructure/token"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
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

	router := mux.NewRouter()

	router.Handle("/images/{picture}", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	posthttpport.HandlePostHttpRoutes(router)
	userhttpport.HandleUserRoutes(router)

	log.Println("Running server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
