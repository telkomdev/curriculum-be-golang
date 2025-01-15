package main

import (
	"50-jwt-authentication/app/adapter/mongodb"
	"50-jwt-authentication/app/controllers"
	"50-jwt-authentication/app/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

// @title           50-JWT-Authentication
// @version         1.0
// @description     Add User Service to manage user, authentication and role authorization
// @host      		localhost:8000
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

	//routing mux
	mux := http.NewServeMux()

	// create new NewUser controllers
	itemHandler := controllers.NewItemsList(mod)

	// create new NewItemsList controllers
	// create new items controllers, add MiddlewareAuth to auth checking
	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/item",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", itemHandler.Get()).   // list all items
			Method("POST", itemHandler.Post()). // add item
			Serve(), mod, []string{"admin", "user"}),
	)

	// create new NewUser controllers
	userHandler := controllers.NewUser(mod)

	// create new UserCreateAdmin controllers
	// add MiddlewareSuperAdminAuth to auth checking
	mux.Handle("/api/v1/user/create/admin",
		controllers.MiddlewareSuperAdminAuth(controllers.NewController(nil).
			Method("POST", userHandler.CreateAdmin()). // find user by id
			Serve()),
	)

	// create new UserCreate controllers
	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user/create",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("POST", userHandler.Create()). // find user by id
			Serve(), mod, []string{"admin"}),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user/me",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", userHandler.Me()). // find user by id
			Serve(), mod, []string{"admin", "user"}),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user/auth",
		controllers.NewController(nil).
			Method("POST", userHandler.Auth()).
			Serve(),
	)

	// add MiddlewareAuth to auth checking
	// userIdPattern - validate match path user
	userIdPattern := regexp.MustCompile(`^/api/v1/user/([A-Za-z0-9]+)$`)
	mux.Handle("/api/v1/user/",
		controllers.MiddlewareAuth(controllers.NewController(userIdPattern).
			Method("GET", userHandler.Get()).       // find user by id
			Method("PUT", userHandler.Put()).       // edit user by id
			Method("DELETE", userHandler.Delete()). // delete user by id
			Serve(), mod, []string{"admin"}),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", userHandler.GetAll()). // get all users
			Serve(), mod, []string{"admin"}),
	)

	// root controllers
	mux.Handle("/", controllers.NewController(nil).
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
