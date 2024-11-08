package api

import (
	"fmt"
	"supplier-backend/db"
	"supplier-backend/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate`
	TillDate   time.Time `json:"tillDate`
	NumPersons int       `json:"numPersons`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (r *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	roomID := c.Params("id")
	//VALIDATE CORRECTNESS OF THE ID
	oid, oidErr := primitive.ObjectIDFromHex(roomID)
	if oidErr != nil {
		return oidErr
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		res := map[string]any{
			"success": false,
			"message": "server error",
		}
		return c.Status(500).JSON(res)
	}

	booking := types.Booking{
		RoomID:     oid,
		UserID:     user.ID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	fmt.Println(booking)

	return nil
}
