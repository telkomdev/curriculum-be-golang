package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type MongoDB struct {
	database *mongo.Database
}

func New() *MongoDB {
	m := &MongoDB{}
	err := m.Init()
	if err != nil {
		log.Panicf("failed to connect mongodb, %s\n", err.Error())
	} // starting mongodb connection
	return m
}

func (c *MongoDB) Init() (err error) {
	// default options
	mongoHost := "localhost"
	mongoPort := "27017"
	mongoDBName := "default_database"
	mongoUsername := "root"
	mongoPassword := "MongoDBPasswordDev"

	// get options form env if available
	if len(os.Getenv("MONGO_HOST")) != 0 {
		mongoHost = os.Getenv("MONGO_HOST")
	}

	if len(os.Getenv("MONGO_PORT")) != 0 {
		mongoPort = os.Getenv("MONGO_PORT")
	}

	if len(os.Getenv("MONGO_DBNAME")) != 0 {
		mongoDBName = os.Getenv("MONGO_DBNAME")
	}

	if len(os.Getenv("MONGO_USERNAME")) != 0 {
		mongoUsername = os.Getenv("MONGO_USERNAME")
	}

	if len(os.Getenv("MONGO_PASSWORD")) != 0 {
		mongoPassword = os.Getenv("MONGO_PASSWORD")
	}

	// parsing and create client mongodb connection
	clientOptions := options.Client()
	clientOptions.Auth = &options.Credential{
		Username: mongoUsername,
		Password: mongoPassword,
	}
	clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort))
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to mongodb
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	c.database = client.Database(mongoDBName)

	return nil
}

func (c *MongoDB) Database() *mongo.Database {
	return c.database
}
