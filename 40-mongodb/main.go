package main

import (
	"40-mongodb/app/adapter/mongodb"
	"40-mongodb/app/controllers"
	"40-mongodb/app/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// @title           40-mongodb
// @version         1.0
// @description     GET item list from database and POST new item (save item data to mongo database)
// @host      		localhost:8000
func main() {
	//define port for serving request
	var address = ":8000"
	//get from ENV if available
	if len(os.Getenv("ADDRESS")) != 0 {
		address = os.Getenv("ADDRESS")
	}

	// create mongodb object
	db := mongodb.New()
	// create models
	mod := models.New(db)

	//routing mux
	mux := http.NewServeMux()

	// create new NewUser controllers
	itemHandler := controllers.NewItemsList(mod)

	// create new NewItemsList controllers
	// create new items controllers, add MiddlewareAuth to auth checking
	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/item",
		controllers.NewController(nil).
			Method("GET", itemHandler.Get()).   // list all items
			Method("POST", itemHandler.Post()). // add item
			Serve(),
	)

	// root controllers
	mux.Handle("/api/v1", controllers.NewController(nil).
		Method("GET", controllers.NewRoot().Get()).
		Serve())

	//start the apps
	server := http.Server{
		Addr:    address,
		Handler: mux,
	}

	//log the response
	log.Printf("server started at %s\n", address)

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}
