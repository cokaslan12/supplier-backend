package db

import (
	"context"
	"supplier-backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	coll *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client, dbName string) *MongoRoomStore {
	return &MongoRoomStore{
		coll: client.Database(dbName).Collection(ROOM_COL),
	}
}

func (r *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := r.coll.InsertOne(ctx, &room)
	if err != nil {
		return nil, err
	}

	room.ID = resp.InsertedID.(primitive.ObjectID)
	filter := bson.M{
		"_id": room.HotelID,
	}

	update := bson.M{
		"$push":bson.M{
			"rooms":rooms.ID,
		},
	}


	return room, nil
}
