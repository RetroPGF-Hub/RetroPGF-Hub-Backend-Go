package usersrepository

import (
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoritePb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	grpcconn "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/grpcConn"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/jwtauth"
	"context"
	"errors"
	"log"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	UsersRepositoryService interface {
		InsertOneUser(pctx context.Context, req *users.UserDb) (primitive.ObjectID, error)
		FindOneUserWithIdWithPassword(pctx context.Context, userId primitive.ObjectID) (*users.UserDb, error)
		IsUniqueUser(pctx context.Context, email string) (bool, error)
		FindOneUserWithEmail(pctx context.Context, email string) (*users.UserDb, error)
		FindOneUserWithId(pctx context.Context, userId primitive.ObjectID) (*users.UserDb, error)
		GetFavProjectByUserId(pctx context.Context, favGrpcUrl, userId string) (*favPb.GetAllFavRes, error)
		FindManyUserId(pctx context.Context, usersId []primitive.ObjectID) ([]*users.UserDb, error)
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

func (r *usersRepository) FindOneUserWithIdWithPassword(pctx context.Context, userId primitive.ObjectID) (*users.UserDb, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(users.UserDb)

	if err := col.FindOne(ctx, bson.M{"_id": userId}).Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("error: id not found")
		}
		return nil, err
	}

	return result, nil
}

func (r *usersRepository) FindOneUserWithId(pctx context.Context, userId primitive.ObjectID) (*users.UserDb, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(users.UserDb)
	projection := bson.M{
		"_id":        1,
		"email":      1,
		"source":     1,
		"profile":    1,
		"user_name":  1,
		"first_name": 1,
		"last_name":  1,
	}

	options := options.FindOne().SetProjection(projection)

	if err := col.FindOne(ctx, bson.M{"_id": userId}, options).Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("error: id not found")
		}
		return nil, err
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

func (r *usersRepository) FindManyUserId(pctx context.Context, usersId []primitive.ObjectID) ([]*users.UserDb, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	matchInState := bson.D{
		{"$match",
			bson.D{
				{"_id",
					bson.D{
						{"$in",
							usersId,
						},
					},
				},
			},
		},
	}
	projectInState := bson.D{
		{"$project",
			bson.D{
				{"_id", 1},
				{"email", 1},
				{"source", 1},
				{"profile", 1},
				{"user_name", 1},
				{"first_name", 1},
				{"last_name", 1},
			},
		},
	}
	pipeline := mongo.Pipeline{matchInState, projectInState}
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rUsers []*users.UserDb

	for cursor.Next(ctx) {
		var result users.UserDb
		err := cursor.Decode(&result)
		if err != nil {
			return nil, errors.New("error: decoding result failed")
		}
		rUsers = append(rUsers, &result)
	}

	return rUsers, nil
}

func (r *usersRepository) GetFavProjectByUserId(pctx context.Context, favGrpcUrl, userId string) (*favPb.GetAllFavRes, error) {
	jwtauth.SetApiKeyInContext(&pctx)
	conn, err := grpcconn.NewGrpcClient(favGrpcUrl)
	if err != nil {
		log.Printf("Error: Grpc Conn Error %s", err.Error())
		return nil, errors.New("error: grpc connection to fav failed")
	}

	result, err := conn.Fav().GetAllFavByUserId(pctx, &favPb.GetAllFavProjectReq{
		UserId: userId,
	})
	if err != nil {
		log.Printf("Error: Create Fav For Project Error %s", err.Error())
		return nil, errors.New("errors")
	}
	return result, nil
}
