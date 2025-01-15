package main

import (
	"30-routing/app/controllers"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// @title           30-routing
// @version         1.0
// @description     Example of routing with golang, adding GET and POST route
// @host      		localhost:8000
func main() {
	//define port for serving request
	var address = ":8000"
	//get from ENV if available
	if len(os.Getenv("ADDRESS")) != 0 {
		address = os.Getenv("ADDRESS")
	}

	//routing mux
	mux := http.NewServeMux()
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
