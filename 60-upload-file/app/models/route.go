package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"
)

// FindRouteByFromAndTo - search route document by from and to
func (c *Models) FindRouteByFromAndTo(collectionName string, from, to string) (Route, error) {
	filter := bson.D{{Key: "from", Value: strings.ToLower(from)}, {Key: "to", Value: strings.ToLower(to)}}
	return c.getSingleRoute(collectionName, filter)
}

// FindRouteById - search route document by id
func (c *Models) FindRouteById(collectionName string, key string) (Route, error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return Route{}, err
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	return c.getSingleRoute(collectionName, filter)
}

func (c *Models) getSingleRoute(collectionName string, filter interface{}) (Route, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res Route
	err := c.mongodb.Database().Collection(collectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	// convert to title
	res.From = cases.Title(language.Und).String(res.From)
	res.To = cases.Title(language.Und).String(res.To)

	return res, nil
}

func (c *Models) UpsetRouteByFromAndTo(collectionName string, data *Route) error {

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	_, err := c.mongodb.Database().Collection(collectionName).
		UpdateOne(ctx, bson.D{{Key: "from", Value: strings.ToLower(data.From)}, {Key: "to", Value: strings.ToLower(data.To)}},
			bson.D{
				{Key: "$setOnInsert", Value: bson.D{
					{Key: "createdAt", Value: data.CreatedAt},
				}},
				{Key: "$set",
					Value: bson.D{
						{Key: "updatedAt", Value: data.UpdatedAt},
						{Key: "price", Value: data.Price},
						{Key: "departureTime", Value: data.DepartureTime},
					}},
			},
			options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (c *Models) GetAllRoutes(collectionName string, from, to string) (res AllRoutes, err error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}
	if len(from) != 0 && len(to) != 0 {
		filter = bson.D{{Key: "from", Value: strings.ToLower(from)}, {Key: "to", Value: strings.ToLower(to)}}
	} else if len(from) != 0 && len(to) == 0 {
		filter = bson.D{{Key: "from", Value: strings.ToLower(from)}}
	} else if len(to) != 0 && len(from) == 0 {
		filter = bson.D{{Key: "to", Value: strings.ToLower(to)}}
	} else {
		filter = bson.M{}
	}

	cursor, err := c.mongodb.Database().Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return res, err
	}

	if err = cursor.All(context.TODO(), &res.Data); err != nil {
		return res, err
	}

	var data []Route
	for _, s := range res.Data {
		s.From = cases.Title(language.Und).String(s.From)
		s.To = cases.Title(language.Und).String(s.To)
		data = append(data, s)
	}

	res.Data = data
	res.Count = len(res.Data)

	return res, nil
}
