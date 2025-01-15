package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"
)

func (c *Models) InsertTicket(data *Ticket) (id string, err error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	res, err := c.mongodb.Database().Collection(TicketCollectionName).InsertOne(ctx, data)
	if err != nil {
		return id, err
	}

	return fmt.Sprintf("%s", res.InsertedID.(primitive.ObjectID).Hex()), nil

}

func (c *Models) UpsetTicketById(data *Ticket) error {
	id, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	_, err = c.mongodb.Database().Collection(TicketCollectionName).
		UpdateOne(ctx, bson.D{{Key: "_id", Value: id}},
			bson.D{
				{Key: "$setOnInsert", Value: bson.D{
					{Key: "createdAt", Value: data.CreatedAt},
				}},
				{Key: "$set",
					Value: bson.D{
						{Key: "updatedAt", Value: data.UpdatedAt},
						{Key: "price", Value: data.Price},
						{Key: "departureTime", Value: data.DepartureTime},
						{Key: "from", Value: strings.ToLower(data.From)},
						{Key: "to", Value: strings.ToLower(data.To)},
						{Key: "userId", Value: data.UserId},
						{Key: "bookingId", Value: data.BookingId},
					}},
			},
			options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// FindTicketById - search ticket document by id
func (c *Models) FindTicketById(key string) (Ticket, error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return Ticket{}, err
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	return c.getSingleTicket(filter)
}

func (c *Models) getSingleTicket(filter interface{}) (Ticket, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res Ticket
	err := c.mongodb.Database().Collection(TicketCollectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	// convert to title
	res.From = cases.Title(language.Und).String(res.From)
	res.To = cases.Title(language.Und).String(res.To)

	return res, nil
}

func (c *Models) DeleteTicketByID(data string) error {
	id, err := primitive.ObjectIDFromHex(data)
	if err != nil {
		return err
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	_, err = c.mongodb.Database().Collection(TicketCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil

}

func (c *Models) GetAllTickets(from, to, userId, bookingId string, pagination PaginationOption) (res AllTickets, err error) {

	filter := bson.M{}
	if len(from) != 0 {
		filter["from"] = bson.M{"$regex": from, "$options": "im"}
	}

	if len(to) != 0 {
		filter["to"] = bson.M{"$regex": to, "$options": "im"}
	}

	if len(userId) != 0 {
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return res, err
		}
		filter["userId"] = bson.M{"$eq": id}
	}

	if len(bookingId) != 0 {
		id, err := primitive.ObjectIDFromHex(bookingId)
		if err != nil {
			return res, err
		}
		filter["bookingId"] = bson.M{"$eq": id}
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := c.mongodb.Database().Collection(TicketCollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return res, err
	}

	skip := (pagination.Page - 1) * pagination.Size
	opts := options.Find()
	opts.SetLimit(int64(pagination.Size)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{"createdAt", 1}})

	cursor, err := c.mongodb.Database().Collection(TicketCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return res, err
	}

	if err = cursor.All(context.TODO(), &res.Data); err != nil {
		return res, err
	}

	var data []Ticket
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
