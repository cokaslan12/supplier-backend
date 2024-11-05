package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT Authing --")

	token := c.Get("X-Api-Token")

	if token == "" {
		res := map[string]any{
			"success": false,
			"errors":  "unauthorized",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(res)
	}

	claims, err := validateToken(token)
	if err != nil {
		return err
	}

	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)

	// check token expiration
	if time.Now().Unix() > expires {
		return fmt.Errorf("token expired")
	}

	return c.Next()
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET") //WE ARE GETTING COMPUTER ENVIRONMENT VARIABLES, OS COMMUNICATES ITH THE OPERATION SYSTEM EXPORT JWT_SECRET="SUPPLIER-BACKEND"
		return []byte(secret), nil
	})

	if err != nil {

		fmt.Println("failed to parse JWT Token: ", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
