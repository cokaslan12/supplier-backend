package db

// MARK: DBURI
const (
	DB_URI      string = "mongodb://localhost:27017"
	DB_TEST_URI string = "mongodb://localhost:27017"
)

// MARK: DBNAME
const (
	DB_NAME      string = "supplier"
	DB_TEST_NAME string = "supplier-test"
)

// MARK: DBCOLLECTION
const (
	USER_COL  string = "users"
	HOTEL_COL string = "hotels"
	ROOM_COL  string = "rooms"
)


type Store struct{
	UserStore UserStore
	HotelStore HotelStore
	RoomStore RoomStore
}
