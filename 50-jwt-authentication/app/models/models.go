package models

import (
	"50-jwt-authentication/app/adapter/mongodb"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// FindItemByID - search items document by id
func (c *Models) FindItemByID(collectionName string, key string) (Item, error) {
	id, _ := primitive.ObjectIDFromHex(key)
	filter := bson.D{{Key: "_id", Value: id}}
	return c.getSingleItem(collectionName, filter)
}

// FindItemByName - search items document by name
func (c *Models) FindItemByName(collectionName string, key string) (Item, error) {
	filter := bson.D{{Key: "name", Value: key}}
	return c.getSingleItem(collectionName, filter)
}

func (c *Models) getSingleItem(collectionName string, filter interface{}) (Item, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res Item
	err := c.mongodb.Database().Collection(collectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Models) GetAllItem(collectionName string) (res AllItem, err error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := c.mongodb.Database().Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return res, err
	}

	if err = cursor.All(context.TODO(), &res.Data); err != nil {
		return res, err
	}

	res.Count = len(res.Data)

	return res, nil
}

func (c *Models) InsertItem(collectionName string, data *Item) error {

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	_, err := c.mongodb.Database().Collection(collectionName).
		UpdateOne(ctx, bson.D{{Key: "name", Value: data.Name}},
			bson.D{
				{Key: "$setOnInsert",
					Value: bson.D{
						{Key: "name", Value: data.Name},
						{Key: "createdAt", Value: data.CreatedAt},
					}},
				{Key: "$set",
					Value: bson.D{
						{Key: "qty", Value: data.Qty},
						{Key: "updatedAt", Value: data.UpdatedAt},
					}},
			},
			options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (c *Models) InsertRole(collectionName string, data *Role) error {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Upset role, if role exist, do not nothing
	_, err := c.mongodb.Database().Collection(collectionName).
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
func (c *Models) FindRoleByName(collectionName string, key string) (Role, error) {
	filter := bson.D{{Key: "name", Value: key}}
	return c.getSingleRole(collectionName, filter)
}

func (c *Models) getSingleRole(collectionName string, filter interface{}) (Role, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res Role
	err := c.mongodb.Database().Collection(collectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

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
