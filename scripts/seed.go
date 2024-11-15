package main

import (
	"context"
	"fmt"
	"log"
	"supplier-backend/api"
	"supplier-backend/db"
	"supplier-backend/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	bookStore  db.BookingStore
	ctx        = context.Background()
)

func seedUser(isAdmin bool, fName, lName, email, password string) *types.User {
	params := types.CreateUser{
		FirstName: fName,
		LastName:  lName,
		Email:     email,
		Password:  password,
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	InsertedUser, InsertedErr := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(InsertedErr)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	return InsertedUser
}

func seedRoom(size string, ss bool, price float64, hotelId primitive.ObjectID) *types.Room {
	room := types.Room{
		Seaside: ss,
		Size:    size,
		Price:   price,
		HotelID: hotelId,
	}

	InsertedRoom, InsertedErr := roomStore.InsertRoom(ctx, &room)
	if InsertedErr != nil {
		log.Fatal(InsertedErr)
	}

	return InsertedRoom
}

func seedBooking(userId primitive.ObjectID, roomId primitive.ObjectID, numPersons int, fromDate time.Time, tillDate time.Time, canceled bool) *types.Booking {
	booking := types.Booking{
		UserID:     userId,
		RoomID:     roomId,
		NumPersons: numPersons,
		FromDate:   fromDate,
		TillDate:   tillDate,
		Canceled:   canceled,
	}

	InsertedBooking, InsertedErr := bookStore.InsertBooking(ctx, &booking)
	if InsertedErr != nil {
		log.Fatal(InsertedErr)
	}

	return InsertedBooking
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, intertedHotelErr := hotelStore.InsertHotel(ctx, &hotel)
	if intertedHotelErr != nil {
		log.Fatal(intertedHotelErr)
	}

	return insertedHotel
}

func main() {
	seedUser(false, "muzaffer", "çokaslan", "cokaslanmuzaffer@gmail.com", "123456")
	admin := seedUser(true, "admin", "admin", "admin@gmail.com", "123456")
	seedHotel("Bellucia", "France", 3)
	seedHotel("The Cozy Hotel", "The Nederlands", 4)
	hotel := seedHotel("Dont Die In Your Sleep", "London", 1)
	room := seedRoom("small", true, 99.9, hotel.ID)
	seedRoom("normal", true, 199.9, hotel.ID)
	seedRoom("kingsize", true, 222.9, hotel.ID)
	seedBooking(admin.ID, room.ID, 1, time.Now(), time.Now().AddDate(0, 0, 2), false)

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

	//MARK: INITIALIZE STORE ""
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookStore = db.NewMongoBookingStore(client)
}
