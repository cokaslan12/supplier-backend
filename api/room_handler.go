package api

import (
	"context"
	"fmt"
	"net/http"
	"supplier-backend/db"
	"supplier-backend/types"
	"supplier-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate`
	TillDate   time.Time `json:"tillDate`
	NumPersons int       `json:"numPersons`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}

	if p.FromDate.Unix() > p.TillDate.Unix() {
		return fmt.Errorf("tillDate cannot be in the past from fromDate")
	}

	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (r *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := r.store.RoomStore.GetRooms(c.Context(), bson.M{})

	if err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    rooms,
	}

	return c.Status(200).JSON(res)
}

func (r *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return params.validate()
	}

	roomID := c.Params("id")
	//VALIDATE CORRECTNESS OF THE ID
	oid, oidErr := primitive.ObjectIDFromHex(roomID)
	if oidErr != nil {
		return oidErr
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return utils.Error{
			Code:    http.StatusInternalServerError,
			Success: false,
			Err:     "server error",
		}
	}

	ok, err := r.isRoomAvailableForBooking(c.Context(), oid, params)

	if err != nil {
		return err
	}

	if !ok {
		res := map[string]any{
			"success": false,
			"message": fmt.Sprintf("room %s already booked", roomID),
		}
		return c.Status(400).JSON(res)
	}

	booking := types.Booking{
		RoomID:     oid,
		UserID:     user.ID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	insertedBooking, err := r.store.BookingStore.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    insertedBooking,
	}

	return c.Status(200).JSON(res)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomId primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID":   roomId,
		"fromDate": bson.M{"$gte": params.FromDate}, //GTE:GREATER THE EXIST
		"tillDate": bson.M{"$lte": params.TillDate}, //LTE:LESSER THE EXIST
	}

	bookings, err := h.store.BookingStore.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}
