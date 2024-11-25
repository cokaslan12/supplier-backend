package api

import (
	"context"
	"log"
	"os"
	"supplier-backend/db"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDb struct {
	client *mongo.Client
	store  *db.Store
}

func (tDb *testDb) tearDown(t *testing.T) {
	dbName := os.Getenv(db.MongoDBNameEnvName)
	if err := tDb.client.Database(dbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() *testDb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}
	mongoEndPoint := os.Getenv("MONGO_DB_TEST_URL")
	client, mongoErr := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(mongoEndPoint))
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
