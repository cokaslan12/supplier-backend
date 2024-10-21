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

	//MARK: STORE INITIALIZATION
	userStore := db.NewMongoUserStore(client)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	store := db.Store{
		UserStore:  userStore,
		HotelStore: hotelStore,
		RoomStore:  roomStore,
	}

	//MARK: HANDLER INITIALIZATION
	userHandler := api.NewUserHandler(&store)
	hotelHandler := api.NewHotelHandler(&store)

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")

	//MARK: USERS API
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id/get", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id/put", userHandler.HandlePutUser)

	//MARK: HOTELS API
	apiV1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id/get", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	err := app.Listen(*listenAddr)

	if err != nil {
		log.Fatal(err)
	}
}
