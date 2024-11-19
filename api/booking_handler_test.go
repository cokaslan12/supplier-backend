package api

import (
	"fmt"
	"supplier-backend/db/fixtures"
	"testing"
	"time"
)

func TestAdminGetBookings(t *testing.T) {
	tDb := setup()
	defer tDb.tearDown(t)

	user := fixtures.AddUser(tDb.store, "muzaffer", "cokaslan", false)
	hotel := fixtures.AddHotel(tDb.store, "Foo Hotel", "Bar Location", 1, nil)
	room := fixtures.AddRoom(tDb.store, true, "large", 30.9, hotel.ID)
	fromDate := time.Now()
	tillDate := time.Now().AddDate(0, 0, 2)
	booking := fixtures.AddBooking(tDb.store, user.ID, room.ID, 1, fromDate, tillDate, false)

	fmt.Println(booking)

}
