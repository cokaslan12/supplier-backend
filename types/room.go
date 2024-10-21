package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType int

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty`
	Seaside bool               `bson:"seaside" json:"seaside"`
	//small, normal,kingsize
	Size    string             `bson:"size" json:"size"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelId" json:"hotelId`
}
