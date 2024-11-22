package middleware

import (
	"fmt"
	"net/http"
	"os"
	"supplier-backend/db"
	"supplier-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {

	return func(c *fiber.Ctx) error {
		fmt.Println("-- JWT Authing")

		token := c.Get("X-Api-Token")

		if token == "" {
			return utils.ErrUnAuthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		// check token expiration
		if time.Now().Unix() > expires {
			return utils.NewError(http.StatusUnauthorized, false, "token expired")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return utils.ErrUnAuthorized()
		}

		// set the current authenticated user to context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}

}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, utils.ErrUnAuthorized()
		}

		secret := os.Getenv("JWT_SECRET") //WE ARE GETTING COMPUTER ENVIRONMENT VARIABLES, OS COMMUNICATES ITH THE OPERATION SYSTEM EXPORT JWT_SECRET="SUPPLIER-BACKEND"
		return []byte(secret), nil
	})

	if err != nil {

		fmt.Println("failed to parse JWT Token: ", err)
		return nil, utils.ErrUnAuthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, utils.ErrUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, utils.ErrUnAuthorized()
	}

	return claims, nil
}
