package middleware

import (
	"fmt"
	"os"

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

	if err := parseToken(token); err != nil {
		return err
	}

	return nil
}

func parseToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET") //WE ARE GETTING COMPUTER ENVIRONMENT VARIABLES, OS COMMUNICATES ITH THE OPERATION SYSTEM
		fmt.Println("NEVER PRINT SECRET: ", secret)
		return []byte(secret), nil
	})

	if err != nil {

		fmt.Println("failed to parse JWT Token: ", err)
		return fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}

	return fmt.Errorf("unauthorized")
}
