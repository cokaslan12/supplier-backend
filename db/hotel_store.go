package db

import (
	"context"
	"errors"
	"supplier-backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	coll *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		coll: client.Database(dbName).Collection(HOTEL_COL),
	}
}

func (h *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := h.coll.InsertOne(ctx, &hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = resp.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (h *MongoHotelStore) Update(ctx context.Context, filter bson.M, params bson.M) error {

	res, err := h.coll.UpdateOne(ctx, filter, params)

	if err != nil {
		return err
	}

	if res.ModifiedCount > 0 {
		return nil
	}

	return errors.New("hotel could not update")
}
