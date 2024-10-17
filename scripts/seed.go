package main

import (
	"context"
	"fmt"
	"log"
	"supplier-backend/db"
	"supplier-backend/types"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//MARK: SET CONTEXT
	ctx := context.TODO()

	//MARK: SETUP MONGO DB
	client, mongoErr := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DB_URI))
	if mongoErr != nil {
		log.Fatal(mongoErr)
	}

	//MARK: INITIALIZE STORE
	hotelStore := db.NewMongoHotelStore(client, db.DB_NAME)
	roomStore := db.NewMongoRoomStore(client, db.DB_NAME)

	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "Narnia",
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 122.9,
		},
	}

	insertedHotel, intertedHotelErr := hotelStore.InsertHotel(ctx, &hotel)
	if intertedHotelErr != nil {
		log.Fatal(intertedHotelErr)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		intertedRoom, insertedRoomErr := roomStore.InsertRoom(ctx, &room)
		if insertedRoomErr != nil {
			log.Fatal(insertedRoomErr)
		}

		fmt.Println(intertedRoom)
	}
}
