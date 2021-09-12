package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saarmornel/reading-list/repo"
	"github.com/saarmornel/reading-list/web"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetReportCaller(true)
	db, err := repo.Init("app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := httprouter.New()
	web.SetupRoutes(router)

	log.Print("db is ready!")
	log.Fatal(http.ListenAndServe(":3001", router))
}
