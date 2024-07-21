package repository

import (
	"context"
	"log"

	"github.com/shiro1n/go-commerce/internal/user/model"
	"github.com/shiro1n/go-commerce/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func init() {
	userCollection = database.GetMongoDBCollection("go-commerce", "users")
}

func CreateUser(ctx context.Context, user model.User) (*mongo.InsertOneResult, error) {
	user.ID = primitive.NewObjectID()
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Could not create user: %v", err)
		return nil, err
	}
	return result, nil
}

func GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Printf("Could not find user: %v", err)
		return user, err
	}
	return user, nil
}

func UpdateUser(ctx context.Context, user model.User) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	result, err := userCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Could not update user: %v", err)
		return nil, err
	}
	return result, nil
}
