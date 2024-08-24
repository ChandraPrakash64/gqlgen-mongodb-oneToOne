package database

import (
	"context"
	"log"

	"example.com/gqlgen/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) CreateUser(name string, city string, pin string) *model.User {
	userCollection := db.client.Database("dbs").Collection("users")
	addressCollection := db.client.Database("dbs").Collection("address")
	generatedUserID := primitive.NewObjectID().Hex()
	generatedAddressID := primitive.NewObjectID().Hex()
	newAddress := model.Address{
		ID:     generatedAddressID,
		City:   &city,
		Pin:    &pin,
		UserID: &generatedUserID,
	}

	newUser := model.User{
		ID:      generatedUserID,
		Name:    &name,
		Address: &newAddress,
	}

	_, err := userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatal(err)
	}

	_, err1 := addressCollection.InsertOne(context.TODO(), newAddress)
	if err1 != nil {
		log.Fatal(err)
	}
	return &newUser
}

func (db *DB) DeleteUser(id string) bool {
	userCollection := db.client.Database("dbs").Collection("users")
	addressCollection := db.client.Database("dbs").Collection("address")
	//idd, _ := primitive.ObjectIDFromHex(id)
	_, err := userCollection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err1 := addressCollection.DeleteOne(context.TODO(), bson.M{"userid": id})
	if err1 != nil {
		log.Fatal(err1)
		return false
	}
	return true
}

func (db *DB) GetUser(id string) *model.User {
	userCollection := db.client.Database("dbs").Collection("users")
	var reqUser model.User
	err := userCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&reqUser)
	if err != nil {
		log.Fatal(err)
	}
	return &reqUser
}

func (db *DB) GetAddressFromUser(id string) *model.Address {
	addressCollection := db.client.Database("dbs").Collection("address")
	var reqAddress model.Address
	err := addressCollection.FindOne(context.TODO(), bson.M{"userid": id}).Decode(&reqAddress)
	if err != nil {
		log.Fatal(err)
	}
	return &reqAddress
}
