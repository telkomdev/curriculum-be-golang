package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// FindUserByEmail - search users document by email
func (c *Models) FindUserByEmail(key string) (UserList, error) {
	filter := bson.D{{Key: "email", Value: key}}
	return c.getSingleUser(filter)
}

// FindUserById - search users document by id
func (c *Models) FindUserById(key string) (UserList, error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return UserList{}, err
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	return c.getSingleUser(filter)
}

func (c *Models) getSingleUser(filter interface{}) (UserList, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res UserList
	err := c.mongodb.Database().Collection(UserCollectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Models) UpsetUserByEmail(data *UserList) error {

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	_, err := c.mongodb.Database().Collection(UserCollectionName).
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

func (c *Models) DeleteUserByID(data string) error {
	id, err := primitive.ObjectIDFromHex(data)
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	_, err = c.mongodb.Database().Collection(UserCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil

}

func (c *Models) GetAllUsers(name string, pagination PaginationOption) (res AllUsers, err error) {

	var filter interface{}
	if len(name) == 0 {
		filter = bson.M{}
	} else {
		filter = bson.D{{Key: "name", Value: bson.M{"$regex": name, "$options": "im"}}}
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := c.mongodb.Database().Collection(UserCollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return res, err
	}

	skip := (pagination.Page - 1) * pagination.Size
	opts := options.Find()
	opts.SetLimit(int64(pagination.Size)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{"createdAt", 1}})

	cursor, err := c.mongodb.Database().Collection(UserCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return res, err
	}

	if err = cursor.All(context.TODO(), &res.Data); err != nil {
		return res, err
	}

	res.CurrentPage = int64(pagination.Page)
	if count < int64(pagination.Size) {
		res.TotalPages = 1
	} else {
		res.TotalPages = count / int64(pagination.Size)
	}
	res.Count = len(res.Data)
	res.TotalCount = count

	return res, nil
}
