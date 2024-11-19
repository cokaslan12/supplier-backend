package main

import (
	"context"
	"fmt"
	"log"
	"supplier-backend/api"
	"supplier-backend/db"
	"supplier-backend/db/fixtures"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//MARK: SET CONTEXT
	ctx := context.TODO()

	//MARK: SETUP MONGO DB
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	//MARK: DROP DATABASE
	if err = client.Database(db.DB_NAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	//MARK: INITIALIZE STORE ""
	store := &db.Store{
		HotelStore:   db.NewMongoHotelStore(client),
		RoomStore:    db.NewMongoRoomStore(client, db.NewMongoHotelStore(client)),
		UserStore:    db.NewMongoUserStore(client),
		BookingStore: db.NewMongoBookingStore(client),
	}
	user := fixtures.AddUser(store, "muzaffer", "cokaslan", false)
	fmt.Println("Muzaffer -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("Admin -> ", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "Bellucia", "France", 5, nil)
	room := fixtures.AddRoom(store, true, "large", 99.9, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, 1, time.Now(), time.Now().AddDate(0, 0, 5), false)
	fmt.Printf("The booking is -> %v", booking)

}
