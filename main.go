package main

import (
	"context"
	"flag"
	"log"
	"supplier-backend/api"
	"supplier-backend/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.Status(400).JSON(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":2000", "The address of backend")
	flag.Parse()

	client, mongoErr := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DB_URI))
	if mongoErr != nil {
		log.Fatal(mongoErr)
	}

	//MARK: HANDLER INITIALIZATION
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, db.DB_NAME))

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")

	//MARK: USERS API
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id/get", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id/put", userHandler.HandlePutUser)

	err := app.Listen(*listenAddr)

	if err != nil {
		log.Fatal(err)
	}
}
