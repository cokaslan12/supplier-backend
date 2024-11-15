package utils

import (
	"fmt"
	"supplier-backend/types"

	"github.com/gofiber/fiber/v2"
)

func AuthenticatedUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return user, nil

}
