package api

import (
	"errors"
	"fmt"
	"supplier-backend/db"
	"supplier-backend/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store db.UserStore
}

func NewUserHandler(store db.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	ctx := c.Context()
	var (
		//values bson.M
		params types.UpdateUser
		userId = c.Params("id")
	)
	//VALIDATE CORRECTNESS OF THE ID
	oid, oidErr := primitive.ObjectIDFromHex(userId)
	if oidErr != nil {
		return oidErr
	}

	filter := bson.M{"_id": oid}
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := h.store.UpdateUser(ctx, filter, params); err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
	}

	return c.Status(200).JSON(res)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := c.Params("id")

	err := h.store.DeleteUser(ctx, userId)

	if err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
	}

	return c.Status(200).JSON(res)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {

	var params types.CreateUser
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if userValidErrs := params.Validate(); len(userValidErrs) > 0 {
		res := map[string]any{
			"success": false,
			"message": "Request failure",
			"errors":  userValidErrs,
		}
		return c.Status(400).JSON(res)
	}

	user, newUserErr := types.NewUserFromParams(params)
	if newUserErr != nil {
		return newUserErr
	}

	user, err := h.store.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    user,
	}
	return c.Status(200).JSON(res)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.store.GetUserById(c.Context(), id)
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
		"data":    user,
	}

	return c.Status(fiber.StatusOK).JSON(res)

}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {

	users, err := h.store.GetUsers(c.Context())
	if err != nil {
		return err
	}

	res := map[string]any{
		"success": true,
		"message": fmt.Sprintf("%v users found", len(users)),
		"data":    users,
	}

	return c.Status(200).JSON(res)
}
