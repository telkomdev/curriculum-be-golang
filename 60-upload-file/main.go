package main

import (
	"60-upload-file/app/adapter/mongodb"
	"60-upload-file/app/models"
	"60-upload-file/app/router"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// @title           			60-upload-file
// @version         			1.0
// @description     			Add Route service to add, import, find and update ticketing route
// @host      					localhost:8000
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
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

	// creating role if not exist
	err := mod.InsertRole("roles", &models.Role{Name: "admin"})
	if err != nil {
		log.Panicf("failed to create user role %s\n", err.Error())
	}

	err = mod.InsertRole("roles", &models.Role{Name: "user"})
	if err != nil {
		log.Panicf("failed to create user role %s\n", err.Error())
	}

	//start the apps
	server := http.Server{
		Addr:    address,
		Handler: router.New(mod),
	}

	//log the response
	log.Printf("server started at %s\n", address)

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}
