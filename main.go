package main

import (
	"context"
	"log"
	"os"
	"supplier-backend/api"
	"supplier-backend/db"
	"supplier-backend/middleware"
	"supplier-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Configuration
//1. MongoDB endpoint
//2. ListenAddress of your HTTP server
//3. JWT secret
//4. MongoDBName

var config = fiber.Config{
	ErrorHandler: utils.ErrorHandler,
}

func main() {
	mongoEndPoint := os.Getenv("MONGO_DB_URL")
	client, mongoErr := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(mongoEndPoint))
	if mongoErr != nil {
		log.Fatal(mongoErr)
	}

	//MARK: STORE INITIALIZATION
	userStore := db.NewMongoUserStore(client)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	bookingStore := db.NewMongoBookingStore(client)
	store := db.Store{
		UserStore:    userStore,
		HotelStore:   hotelStore,
		RoomStore:    roomStore,
		BookingStore: bookingStore,
	}

	//MARK: HANDLER INITIALIZATION
	authHandler := api.NewAuthHandler(&store)
	userHandler := api.NewUserHandler(&store)
	hotelHandler := api.NewHotelHandler(&store)
	roomHandler := api.NewRoomHandler(&store)
	bookingHandler := api.NewBookingHandler(&store)

	app := fiber.New(config)
	auth := app.Group("/api")
	apiV1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	admin := apiV1.Group("/admin", middleware.AdminAuth)

	//MARK: AUTH API
	auth.Post("/auth", authHandler.HandleAuthenticate)

	//VERSIONED API ROUTES
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

	//MARK: ROOMS API
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiV1.Get("/rooms", roomHandler.HandleGetRooms)

	//MARK: BOOKING API
	//ADMIN
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	//USERS
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")

	err := app.Listen(listenAddr)

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
