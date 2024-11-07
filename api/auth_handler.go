package api

import (
	"context"
	"errors"
	"fmt"
	"os"
	"supplier-backend/db"
	"supplier-backend/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
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

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func invalidCredentials(c *fiber.Ctx) error {
	res := map[string]any{
		"success": false,
		"error":   "invalid credentials",
	}

	return c.Status(400).JSON(res)
}

// A handler should do:
// - serialization of the incoming request (JSON)
// - do some data fetching from db
// - call some business logic
// - return the data back the user
func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	ctx := context.Background()
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return invalidCredentials(c)
	}

	user, err := h.store.UserStore.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}

	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	res := map[string]any{
		"success": true,
		"message": "Request successfuly",
		"data":    resp,
	}

	return c.Status(200).JSON(res)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	fmt.Println("---- ", secret, " ----")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret")

	}
	return tokenStr
}
