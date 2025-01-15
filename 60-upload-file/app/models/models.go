package models

import (
	"60-upload-file/app/adapter/mongodb"
	"os"
)

// Secret - default secret
var Secret string = "aOFNMxyVIZfAANsT"

// SecretKey - default secret key
var SecretKey string = "BctbLulGvxijNQKi"

func init() {
	//get from ENV if available
	if len(os.Getenv("SECRET")) != 0 {
		Secret = os.Getenv("SECRET")
	}

	if len(os.Getenv("SECRET_KEY")) != 0 {
		SecretKey = os.Getenv("SECRET_KEY")
	}
}

type Models struct {
	mongodb *mongodb.MongoDB // mongodb object
}

// New - return new mongodb model object
func New(mongodb *mongodb.MongoDB) *Models {
	return &Models{mongodb: mongodb}
}
