package models

import (
	"100-ticketing/app/adapter/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
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

func (c *Models) CreateAllIndex() (err error) {
	var indexModel = make(map[string][]mongo.IndexModel)
	for _, s := range IndexModels {
		optAsc := options.Index()
		optAsc.SetUnique(s.Unique)

		optsDesc := options.Index()
		optsDesc.SetUnique(s.Unique)

		indexModel[s.Collection] = append(indexModel[s.Collection], mongo.IndexModel{
			Keys:    bson.M{s.Field: 1},
			Options: optAsc,
		}, mongo.IndexModel{
			Keys:    bson.M{s.Field: -1},
			Options: optsDesc,
		})
	}

	for k, s := range indexModel {
		if _, err := c.createIndex(k, s); err != nil {
			return err
		}
	}

	return err
}

func (c *Models) createIndex(collectionName string, indexModels []mongo.IndexModel) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := c.mongodb.Database().Collection(collectionName)
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return false, err
	}
	return true, nil
}
