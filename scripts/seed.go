package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"supplier-backend/api"
	"supplier-backend/db"
	"supplier-backend/db/fixtures"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	//MARK: SET CONSTANT VALUES
	var (
		ctx           = context.TODO()
		mongoEndPoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)

	//MARK: SETUP MONGO DB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndPoint))
	if err != nil {
		log.Fatal(err)
	}

	//MARK: DROP DATABASE
	if err = client.Database(mongoDBName).Drop(ctx); err != nil {
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

	for i := 0; i < 300; i++ {
		fixtures.AddHotel(store, fmt.Sprintf("Bellucia %v", i), fmt.Sprintf("France %v", i), 5, nil)
	}

}
