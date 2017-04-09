package main

import (
	"log"
	"net/http"

	"github.com/alvinantonius/go-microservice-sample/cmd/contactapp/handler"
	// "github.com/alvinantonius/go-microservice-sample/internal/cache"
	"github.com/alvinantonius/go-microservice-sample/internal/config"
	"github.com/alvinantonius/go-microservice-sample/internal/contacts"
	"github.com/alvinantonius/go-microservice-sample/internal/database"

	"github.com/julienschmidt/httprouter"
)

func init() {
	// init read config
	conf := config.ReadConfig(
		"/etc/config/config-development.json",
		"../../files/config/config-development.json",
	)

	// open database connection
	database.ConnectDB(conf.Database)
}

func main() {

	conf := config.Get()

	handler.Init()

	// init contacts package
	// why i use InitContacts instead of init() ?
	// i'll explain later
	contacts.New()

	// open redis connection
	// cache.ConnectRedis(conf.Redis)

	// router obj
	router := httprouter.New()
	router.POST("/v1/contacts", handler.NewContact)
	router.PATCH("/v1/contacts/:contact_id", handler.UpdateContact)
	router.GET("/v1/contacts", handler.ListContact)
	router.GET("/v1/contacts/:contact_id", handler.GetContact)
	router.DELETE("/v1/contacts/:contact_id", handler.DeleteContact)

	//run http server
	log.Fatal(http.ListenAndServe(conf.Port, router))
}
