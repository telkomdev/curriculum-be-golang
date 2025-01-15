package main

import (
	"100-ticketing/app/adapter/mongodb"
	"100-ticketing/app/models"
	"100-ticketing/app/router"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// @title           			Ticketing App API
// @version         			1.0
// @description     			This ticketing API have feature to manage user, manage route, manage booking. This application developed by using stack Golang, JWT, MongoDB.
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

	// create mongo index for sort and faster query
	err := mod.CreateAllIndex()
	if err != nil {
		log.Panicf("failed to create index %s\n", err.Error())
	}

	// creating role if not exist
	err = mod.InsertRole(&models.Role{Name: "admin"})
	if err != nil {
		log.Panicf("failed to create user role %s\n", err.Error())
	}

	err = mod.InsertRole(&models.Role{Name: "user"})
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
