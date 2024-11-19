package api

import (
	"context"
	"log"
	"supplier-backend/db"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDb struct {
	client *mongo.Client
	store  *db.Store
}

func (tDb *testDb) tearDown(t *testing.T) {
	if err := tDb.client.Database(db.DB_NAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() *testDb {
	client, mongoErr := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DB_TEST_URI))
	if mongoErr != nil {
		log.Fatal(mongoErr)
	}

	return &testDb{
		client: client,
		store: &db.Store{
			HotelStore:   db.NewMongoHotelStore(client),
			RoomStore:    db.NewMongoRoomStore(client, db.NewMongoHotelStore(client)),
			UserStore:    db.NewMongoUserStore(client),
			BookingStore: db.NewMongoBookingStore(client),
		},
	}

}
