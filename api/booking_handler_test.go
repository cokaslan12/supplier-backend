package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"supplier-backend/db/fixtures"
	"supplier-backend/middleware"
	"supplier-backend/types"
	"supplier-backend/utils"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	tDb := setup()
	defer tDb.tearDown(t)
	var (
		nonAuthUser     = fixtures.AddUser(tDb.store, "foo", "bar", false)
		user            = fixtures.AddUser(tDb.store, "muzaffer", "cokaslan", false)
		hotel           = fixtures.AddHotel(tDb.store, "Foo Hotel", "Bar Location", 1, nil)
		room            = fixtures.AddRoom(tDb.store, true, "large", 30.9, hotel.ID)
		fromDate        = time.Now()
		tillDate        = time.Now().AddDate(0, 0, 2)
		insertedBooking = fixtures.AddBooking(tDb.store, user.ID, room.ID, 1, fromDate, tillDate, false)
		app             = fiber.New(fiber.Config{ErrorHandler: utils.ErrorHandler})
		apiRoute        = app.Group("/", middleware.JWTAuthentication(tDb.store.UserStore))
		bookingHandler  = NewBookingHandler(tDb.store)
	)

	apiRoute.Get("/:id", bookingHandler.HandleGetBooking)
	reg := httptest.NewRequest("GET", fmt.Sprintf("/%s", insertedBooking.ID.Hex()), nil)
	reg.Header.Add("Content-Type", "application/json")
	reg.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, regErr := app.Test(reg)
	if regErr != nil {
		t.Error(regErr)
	}

	if resp.StatusCode != 200 {
		t.Error("status code is not ok", resp.StatusCode)
	}

	bodyBytes, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		t.Error(ioErr)
	}

	var m map[string]any
	unMarshallErr := json.Unmarshal(bodyBytes, &m)
	if unMarshallErr != nil {
		t.Error(unMarshallErr)
	}

	bookingData, exist := m["data"]
	if !exist {
		t.Error("Bookings not exist")
	}

	bookingJson, err := json.Marshal(bookingData)
	if err != nil {
		t.Error("Error marshalling map to JSON:", err)
	}

	var booking *types.Booking
	if err := json.Unmarshal(bookingJson, &booking); err != nil {
		t.Error("Error unmarshalling JSON to struct:", err)
	}

	if insertedBooking.ID != booking.ID {
		t.Fatalf("expected %s got %s", insertedBooking.ID, booking.ID)
	}

	//test non auth user cannot access the bookings
	reg = httptest.NewRequest("GET", fmt.Sprintf("/%s", insertedBooking.ID.Hex()), nil)
	reg.Header.Add("Content-Type", "application/json")
	reg.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))

	resp, regErr = app.Test(reg)
	if regErr != nil {
		t.Error(regErr)
	}

	if resp.StatusCode == 200 {
		t.Fatalf("expected a non status code got %d", resp.StatusCode)
	}

}

func TestAdminGetBookings(t *testing.T) {
	tDb := setup()
	defer tDb.tearDown(t)
	var (
		adminUser      = fixtures.AddUser(tDb.store, "admin", "admin", true)
		user           = fixtures.AddUser(tDb.store, "muzaffer", "cokaslan", false)
		hotel          = fixtures.AddHotel(tDb.store, "Foo Hotel", "Bar Location", 1, nil)
		room           = fixtures.AddRoom(tDb.store, true, "large", 30.9, hotel.ID)
		fromDate       = time.Now()
		tillDate       = time.Now().AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(tDb.store, user.ID, room.ID, 1, fromDate, tillDate, false)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tDb.store.UserStore), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tDb.store)
	)

	admin.Get("/bookings", bookingHandler.HandleGetBookings)
	reg := httptest.NewRequest("GET", "/bookings", nil)
	reg.Header.Add("Content-Type", "application/json")
	reg.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, regErr := app.Test(reg)
	if regErr != nil {
		t.Error(regErr)
	}

	if resp.StatusCode != 200 {
		t.Error("status code is not ok", resp.StatusCode)
	}

	bodyBytes, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		t.Error(ioErr)
	}

	var m map[string]any
	unMarshallErr := json.Unmarshal(bodyBytes, &m)
	if unMarshallErr != nil {
		t.Error(unMarshallErr)
	}

	bookingsData, exist := m["data"]
	if !exist {
		t.Error("Bookings not exist")
	}

	bookingsJson, err := json.Marshal(bookingsData)
	if err != nil {
		t.Error("Error marshalling map to JSON:", err)
	}

	var bookings []*types.Booking
	if err := json.Unmarshal(bookingsJson, &bookings); err != nil {
		t.Error("Error unmarshalling JSON to struct:", err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s but got %s", booking.ID, have.ID)
	}

	//test non admin cannot access the bookings
	reg = httptest.NewRequest("GET", "/bookings", nil)
	reg.Header.Add("Content-Type", "application/json")
	reg.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, regErr = app.Test(reg)
	if regErr != nil {
		t.Error(regErr)
	}

	if resp.StatusCode == 200 {
		t.Fatalf("expected a non status code got %d", resp.StatusCode)
	}

}
