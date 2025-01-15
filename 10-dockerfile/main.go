package main

//import module/library
import (
	"fmt"
	"io"
	"net/http"
)

// define request handler and response
func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	//define port for serving request
	var address = ":8080"

	//log the response
	fmt.Printf("server started at %s\n", address)

	//handler for request
	http.HandleFunc("/", helloHandler)

	//start the apps
	http.ListenAndServe(address, nil)
}
