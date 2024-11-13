package api

import (
	"supplier-backend/db"
	"supplier-backend/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// TODO: this needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.BookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    bookings,
	}

	return c.Status(200).JSON(res)
}

// TODO: this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.BookingStore.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    booking,
	}

	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return err
	}

	if booking.UserID != user.ID {
		res := map[string]any{
			"success": true,
			"message": "not authorized",
			"data":    booking,
		}
		return c.Status(401).JSON(res)
	}

	return c.Status(200).JSON(res)
}
