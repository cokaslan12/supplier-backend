package db

import (
	"context"
	"os"
	"supplier-backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	coll       *mongo.Collection
	HotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	dbName := os.Getenv(MongoDBNameEnvName)
	return &MongoRoomStore{
		coll:       client.Database(dbName).Collection(ROOM_COL),
		HotelStore: hotelStore,
	}
}

func (r *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := r.coll.InsertOne(ctx, &room)
	if err != nil {
		return nil, err
	}

	room.ID = resp.InsertedID.(primitive.ObjectID)
	filter := Map{"_id": room.HotelID}
	update := Map{"$push": bson.M{"rooms": room.ID}}

	if err := r.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}

func (r *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room

	if curErr := cur.All(ctx, &rooms); curErr != nil {
		return nil, curErr
	}

	return rooms, nil

}
