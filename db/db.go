package db

const MongoDBNameEnvName = "MONGO_DB_NAME"

// MARK: DBCOLLECTION
const (
	USER_COL    string = "users"
	HOTEL_COL   string = "hotels"
	ROOM_COL    string = "rooms"
	BOOKING_COL string = "bookings"
)

type Store struct {
	UserStore    UserStore
	HotelStore   HotelStore
	RoomStore    RoomStore
	BookingStore BookingStore
}
