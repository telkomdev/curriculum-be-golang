package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (c *Models) InsertRole(data *Role) error {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Upset role, if role exist, do not nothing
	_, err := c.mongodb.Database().Collection(RoleCollectionName).
		UpdateOne(ctx, bson.D{{Key: "name", Value: data.Name}},
			bson.D{
				{Key: "$setOnInsert",
					Value: bson.D{
						{Key: "name", Value: data.Name},
					}},
			},
			options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// FindRoleByName - search roles document by name
func (c *Models) FindRoleByName(key string) (Role, error) {
	filter := bson.D{{Key: "name", Value: key}}
	return c.getSingleRole(filter)
}

func (c *Models) getSingleRole(filter interface{}) (Role, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res Role
	err := c.mongodb.Database().Collection(RoleCollectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
