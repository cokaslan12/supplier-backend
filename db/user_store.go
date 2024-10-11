package db

import (
	"context"
	"errors"

	"supplier-backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(ctx context.Context, filter bson.M, update types.UpdateUser) error
}

type MongoUserStore struct {
	coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		coll: client.Database(DB_NAME).Collection(USER_COL),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	//VALIDATE CORRECTNESS OF THE ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {

	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if curErr := cur.All(ctx, &users); curErr != nil {
		return nil, curErr
	}

	return users, nil

}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	//VALIDATE CORRECTNESS OF THE ID
	oid, oidErr := primitive.ObjectIDFromHex(id)
	if oidErr != nil {
		return oidErr
	}

	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount > 0 {
		return nil
	}

	return errors.New("user could not delete")
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, values types.UpdateUser) error {

	updates := bson.D{
		{Key: "$set", Value: values.ToBSON()},
	}
	res, err := s.coll.UpdateOne(ctx, filter, updates)

	if err != nil {
		return err
	}

	if res.ModifiedCount > 0 {
		return nil
	}

	return errors.New("user has been update")
}
