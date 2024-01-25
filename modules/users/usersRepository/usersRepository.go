package usersrepository

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UsersRepositoryService interface {
		InsertOneUser(pctx context.Context, req *users.UserDb) (primitive.ObjectID, error)
		FindOneUserWithId(pctx context.Context, userId primitive.ObjectID) (*users.UserDb, error)
		IsUniqueUser(pctx context.Context, email string) (bool, error)
		FindOneUserWithEmail(pctx context.Context, email string) (*users.UserDb, error)
	}

	usersRepository struct {
		db *mongo.Client
	}
)

func NewUsersRepository(db *mongo.Client) UsersRepositoryService {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) userDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("user_db")
}

func (r *usersRepository) InsertOneUser(pctx context.Context, req *users.UserDb) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.userDbConn(ctx)
	col := db.Collection("users")

	userId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneUser: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one user fail")
	}

	return userId.InsertedID.(primitive.ObjectID), nil
}

func (r *usersRepository) FindOneUserWithId(pctx context.Context, userId primitive.ObjectID) (*users.UserDb, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(users.UserDb)

	if err := col.FindOne(ctx, bson.M{"_id": userId}).Decode(result); err != nil {
		return nil, errors.New("error: id not found")
	}

	return result, nil
}

func (r *usersRepository) FindOneUserWithEmail(pctx context.Context, email string) (*users.UserDb, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(users.UserDb)

	if err := col.FindOne(ctx, bson.M{"email": email}).Decode(result); err != nil {
		return nil, errors.New("error: id not found")
	}

	return result, nil
}

func (r *usersRepository) IsUniqueUser(pctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	filter := bson.M{"email": email}
	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.New("error: checking user unique failed")
	}
	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}

}
