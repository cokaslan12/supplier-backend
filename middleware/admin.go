package middleware

import (
	"supplier-backend/types"
	"supplier-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return utils.ErrUnAuthorized()
	}

	if !user.IsAdmin {
		return utils.ErrUnAuthorized()
	}

	return c.Next()
}
