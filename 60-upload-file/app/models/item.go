package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

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
