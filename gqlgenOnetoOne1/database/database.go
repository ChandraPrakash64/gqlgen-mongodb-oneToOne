package database

import (
	"context"
	"log"

	"example.com/aadahar/graph/model"
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
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) CreateUser(input model.CreateInput) *model.User {
	userCollection := db.client.Database("familyTree").Collection("users")
	familyCollection := db.client.Database("familyTree").Collection("family")
	generatedFamilyId := primitive.NewObjectID().Hex()
	generatedUserId := primitive.NewObjectID().Hex()
	newFamily := model.Family{
		ID:            &generatedFamilyId,
		TotalMembers:  input.TotalMembers,
		MaleMembers:   input.MaleMembers,
		FemaleMembers: input.FemaleMembers,
		UserID:        &generatedUserId,
	}
	newUser := model.User{
		ID:     &generatedUserId,
		Name:   input.Name,
		Family: &newFamily,
	}
	_, err := userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatal(err)
	}
	_, err1 := familyCollection.InsertOne(context.TODO(), newFamily)
	if err1 != nil {
		log.Fatal(err1)
	}
	return &newUser
}

/*func (db *DB) updateUser(input model.CreateInput) *model.User {
	userCollection := db.client.Database("familyTree").Collection("users")
	familyCollection := db.client.Database("familyTree").Collection("family")
	var updatedFamily model.Family
	if input.FemaleMembers != nil {
		updatedFamily.FemaleMembers = input.FemaleMembers
	}
	if input.MaleMembers != nil {
		updatedFamily.MaleMembers = input.MaleMembers
	}
	if input.TotalMembers != nil {
		updatedFamily.TotalMembers = input.TotalMembers
	}

} */

func (db *DB) DeleteUser(id string) bool {
	userCollection := db.client.Database("familyTree").Collection("users")
	familyCollection := db.client.Database("familyTree").Collection("family")
	_, err := userCollection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err2 := familyCollection.DeleteOne(context.TODO(), bson.M{"userid": id})
	if err2 != nil {
		log.Fatal(err)
		return false
	}
	return true

}

func (db *DB) GetAllUsers() []*model.User {
	userCollection := db.client.Database("familyTree").Collection("users")
	var listUsers []*model.User
	cursor, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	cursor.All(context.TODO(), &listUsers)
	return listUsers
}

func (db *DB) GetUser(id string) *model.User {
	userCollection := db.client.Database("familyTree").Collection("users")
	var reqUser model.User
	err := userCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&reqUser)
	if err != nil {
		log.Fatal(err)
	}
	return &reqUser
}

func (db *DB) GetFamily(id string) *model.Family {
	familyCollection := db.client.Database("familyTree").Collection("family")
	var reqFamily model.Family
	err := familyCollection.FindOne(context.TODO(), bson.M{"userid": id}).Decode(&reqFamily)
	if err != nil {
		log.Fatal(err)
	}
	return &reqFamily
}
