package api

import (
	"supplier-backend/db"
	"supplier-backend/utils"

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

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.BookingStore.GetBookingByID(c.Context(), id)
	if err != nil {
		return utils.ErrResourceNotFound("booking")
	}

	user, err := utils.AuthenticatedUser(c)

	if err != nil {
		return utils.ErrUnAuthorized()
	}

	if booking.UserID != user.ID {
		return utils.ErrUnAuthorized()
	}

	if err := h.store.BookingStore.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
	}

	return c.Status(200).JSON(res)
}

// TODO: this needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.BookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return utils.ErrResourceNotFound("bookings")
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
		return utils.ErrResourceNotFound("booking")
	}
	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    booking,
	}

	user, err := utils.AuthenticatedUser(c)

	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return utils.ErrUnAuthorized()
	}

	return c.Status(200).JSON(res)
}
