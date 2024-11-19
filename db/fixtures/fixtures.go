package fixtures

import (
	"context"
	"fmt"
	"log"
	"supplier-backend/db"
	"supplier-backend/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, userId, roomId primitive.ObjectID, numPersons int, fromDate, tillDate time.Time, canceled bool) *types.Booking {
	booking := types.Booking{
		UserID:     userId,
		RoomID:     roomId,
		NumPersons: numPersons,
		FromDate:   fromDate,
		TillDate:   tillDate,
		Canceled:   canceled,
	}

	InsertedBooking, InsertedErr := store.BookingStore.InsertBooking(context.TODO(), &booking)
	if InsertedErr != nil {
		log.Fatal(InsertedErr)
	}

	return InsertedBooking
}

func AddRoom(store *db.Store, ss bool, size string, price float64, hid primitive.ObjectID) *types.Room {
	room := types.Room{
		Seaside: ss,
		Size:    size,
		Price:   price,
		HotelID: hid,
	}

	InsertedRoom, InsertedErr := store.RoomStore.InsertRoom(context.TODO(), &room)
	if InsertedErr != nil {
		log.Fatal(InsertedErr)
	}

	return InsertedRoom
}

func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIds = rooms
	if roomIds == nil {
		roomIds = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIds,
		Rating:   rating,
	}

	insertedHotel, intertedHotelErr := store.HotelStore.InsertHotel(context.TODO(), &hotel)
	if intertedHotelErr != nil {
		log.Fatal(intertedHotelErr)
	}

	return insertedHotel
}

func AddUser(store *db.Store, fName, lName string, isAdmin bool) *types.User {
	params := types.CreateUser{
		FirstName: fName,
		LastName:  lName,
		Email:     fmt.Sprintf("%s@%s.com", fName, lName),
		Password:  fmt.Sprintf("%s_%s", fName, lName),
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	InsertedUser, InsertedErr := store.UserStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(InsertedErr)
	}
	return InsertedUser
}
