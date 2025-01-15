package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (c *Models) InsertBooking(data *Booking) (id string, err error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// inserting or updating document
	data.CreatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")

	res, err := c.mongodb.Database().Collection(BookingCollectionName).InsertOne(ctx, data)
	if err != nil {
		return id, err
	}

	return fmt.Sprintf("%s", res.InsertedID.(primitive.ObjectID).Hex()), nil

}

// FindBookingById - search bookings document by id
func (c *Models) FindBookingById(key string) (Booking, error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return Booking{}, err
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$eq": id}}}
	return c.getSingleBooking(filter)
}

func (c *Models) getSingleBooking(filter interface{}) (Booking, error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get single result data
	var res Booking
	err := c.mongodb.Database().Collection(BookingCollectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Models) UpdateBookingById(data *Booking) (err error) {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return err
	}

	data.UpdatedAt = time.Now().Format("2006-01-02T15:04:05.000Z")
	_, err = c.mongodb.Database().Collection(BookingCollectionName).UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "tickets", Value: data.Tickets},
				{Key: "paymentStatus", Value: data.PaymentStatus},
				{Key: "updatedAt", Value: data.UpdatedAt},
			}},
	})
	if err != nil {
		return err
	}

	return nil

}

func (c *Models) GetAllBookings(pagination PaginationOption) (res AllBookings, err error) {

	filter := bson.M{} // create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := c.mongodb.Database().Collection(BookingCollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return res, err
	}

	skip := (pagination.Page - 1) * pagination.Size
	opts := options.Find()
	opts.SetLimit(int64(pagination.Size)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{"createdAt", 1}})

	cursor, err := c.mongodb.Database().Collection(BookingCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return res, err
	}

	if err = cursor.All(context.TODO(), &res.Data); err != nil {
		return res, err
	}

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
