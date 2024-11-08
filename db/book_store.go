package db

import "go.mongodb.org/mongo-driver/mongo"

type BookStore interface {
	//
}

type MongoBookStore struct {
	coll *mongo.Collection
}

func NewMongoBookStore(client *mongo.Client) *MongoBookStore {
	return &MongoBookStore{
		coll: client.Database(DB_URI).Collection(BOOK_COL),
	}
}
