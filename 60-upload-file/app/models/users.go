package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// FindUserByEmail - search users document by email
func (c *Models) FindUserByEmail(collectionName string, key string) (UserList, error) {
	filter := bson.D{{Key: "email", Value: key}}
	return c.getSingleUser(collectionName, filter)
}

// FindUserById - search users document by id
func (c *Models) FindUserById(collectionName string, key string) (UserList, error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return UserList{}, err
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	return c.getSingleUser(collectionName, filter)
}

func (c *Models) getSingleUser(collectionName string, filter interface{}) (UserList, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res UserList
	err := c.mongodb.Database().Collection(collectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Models) UpsetUserByEmail(collectionName string, data *UserList) error {

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	_, err := c.mongodb.Database().Collection(collectionName).
		UpdateOne(ctx, bson.D{{Key: "email", Value: data.Email}},
			bson.D{
				{Key: "$setOnInsert", Value: bson.D{
					{Key: "password", Value: data.Password},
					{Key: "roles", Value: data.Roles},
					{Key: "createdAt", Value: data.CreatedAt},
				}},
				{Key: "$set",
					Value: bson.D{
						{Key: "name", Value: data.Name},
						{Key: "updatedAt", Value: data.UpdatedAt},
					}},
			},
			options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (c *Models) DeleteUserByID(collectionName string, data string) error {
	id, err := primitive.ObjectIDFromHex(data)
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	_, err = c.mongodb.Database().Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil

}

func (c *Models) GetAllUsers(collectionName string, name string) (res AllUsers, err error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}
	if len(name) == 0 {
		filter = bson.M{}
	} else {
		filter = bson.D{{Key: "name", Value: name}}
	}

	cursor, err := c.mongodb.Database().Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return res, err
	}

	if err = cursor.All(context.TODO(), &res.Data); err != nil {
		return res, err
	}

	res.Count = len(res.Data)

	return res, nil
}
