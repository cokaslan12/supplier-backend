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
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	GetHotelById(context.Context, string) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	coll *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		coll: client.Database(DB_NAME).Collection(HOTEL_COL),
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

func (h *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {

	cur, err := h.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if curErr := cur.All(ctx, &hotels); curErr != nil {
		return nil, curErr
	}

	return hotels, nil
}

func (h *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	//VALIDATE CORRECTNESS OF THE ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := h.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
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
