package api

import (
	"errors"
	"fmt"
	"supplier-backend/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	//FILTER
	filter := bson.M{}

	//GET DATA
	hotels, err := h.store.HotelStore.GetHotels(c.Context(), filter)
	if err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": fmt.Sprintf("%v hotels found", len(hotels)),
		"data":    hotels,
	}

	return c.Status(200).JSON(res)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.HotelStore.GetHotelById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			res := map[string]any{
				"success": false,
				"error":   "not found",
			}

			return c.Status(400).JSON(res)
		}
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    hotel,
	}

	return c.Status(fiber.StatusOK).JSON(res)

}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {

	id := c.Params("id")
	//VALIDATE CORRECTNESS OF THE ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotelId": oid}

	rooms, err := h.store.RoomStore.GetRooms(c.Context(), filter)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			res := map[string]any{
				"success": false,
				"error":   "not found",
			}

			return c.Status(400).JSON(res)
		}
		return err
	}

	res := map[string]any{
		"success": true,
		"message": fmt.Sprintf("%v rooms found", len(rooms)),
		"data":    rooms,
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
