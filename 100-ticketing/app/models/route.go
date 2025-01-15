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
func (c *Models) FindRouteByFromAndTo(from, to string) (Route, error) {
	filter := bson.D{{Key: "from", Value: strings.ToLower(from)}, {Key: "to", Value: strings.ToLower(to)}}
	return c.getSingleRoute(filter)
}

// FindRouteById - search route document by id
func (c *Models) FindRouteById(key string) (Route, error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return Route{}, err
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	return c.getSingleRoute(filter)
}

func (c *Models) getSingleRoute(filter interface{}) (Route, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res Route
	err := c.mongodb.Database().Collection(RouteCollectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	// convert to title
	res.From = cases.Title(language.Und).String(res.From)
	res.To = cases.Title(language.Und).String(res.To)

	return res, nil
}

func (c *Models) UpsetRouteByFromAndTo(data *Route) error {

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	_, err := c.mongodb.Database().Collection(RouteCollectionName).
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

func (c *Models) GetAllRoutes(from, to string, pagination PaginationOption) (res AllRoutes, err error) {

	filter := bson.M{}
	if len(from) != 0 {
		filter["from"] = bson.M{"$regex": from, "$options": "im"}
	}

	if len(to) != 0 {
		filter["to"] = bson.M{"$regex": to, "$options": "im"}
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := c.mongodb.Database().Collection(RouteCollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return res, err
	}

	skip := (pagination.Page - 1) * pagination.Size
	opts := options.Find()
	opts.SetLimit(int64(pagination.Size)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{"createdAt", 1}})

	cursor, err := c.mongodb.Database().Collection(RouteCollectionName).Find(ctx, filter, opts)
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
	res.TotalCount = count
	res.CurrentPage = int64(pagination.Page)
	if count < int64(pagination.Size) {
		res.TotalPages = 1
	} else {
		res.TotalPages = count / int64(pagination.Size)
	}

	return res, nil
}
