package main

import (
	"context"
	"log"
	"supplier-backend/db"
	"supplier-backend/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 199.9,
		},
		{
			Size:  "kingsize",
			Price: 222.9,
		},
	}

	insertedHotel, intertedHotelErr := hotelStore.InsertHotel(ctx, &hotel)
	if intertedHotelErr != nil {
		log.Fatal(intertedHotelErr)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, insertedRoomErr := roomStore.InsertRoom(ctx, &room)
		if insertedRoomErr != nil {
			log.Fatal(insertedRoomErr)
		}

	}

}

func main() {
	seedHotel("Bellucia", "France", 3)
	seedHotel("The Cozy Hotel", "The Nederlands", 4)
	seedHotel("Dont Die In Your Sleep", "London", 1)
}

func init() {
	//MARK: SET CONTEXT
	ctx := context.TODO()

	var err error
	//MARK: SETUP MONGO DB
	client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	//MARK: DROP DATABASE
	if err = client.Database(db.DB_NAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	//MARK: INITIALIZE STORE
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
