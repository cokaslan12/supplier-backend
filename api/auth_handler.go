package api

import (
	"context"
	"errors"
	"fmt"
	"supplier-backend/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store *db.Store
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	ctx := context.Background()
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			res := map[string]any{
				"success": false,
				"error":   "invalid credentials",
			}

			return c.Status(400).JSON(res)
		}
		return err
	}

	user, err := h.store.UserStore.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return err
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
	if compareErr != nil {
		return fmt.Errorf("invalid credentials")
	}

	fmt.Println("authenticated -> ", user)

	return nil

}
