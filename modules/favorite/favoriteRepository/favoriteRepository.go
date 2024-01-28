package favoriterepository

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	FavoriteRepositoryService interface {
		PushUserToFav(pctx context.Context, projectId primitive.ObjectID, userId string) (string, error)
		PullUserToFav(pctx context.Context, projectId primitive.ObjectID, userId string) (string, error)
		CountFav(pctx context.Context, projectId primitive.ObjectID, userId string) (int64, int64, error)
		InsertOneFav(pctx context.Context, req *favorite.FavModel) error
	}

	favoriteRepository struct {
		db *mongo.Client
	}
)

func (r *favoriteRepository) favoriteDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("favorite_db")
}

func NewFavoriteRepository(db *mongo.Client) FavoriteRepositoryService {
	return &favoriteRepository{
		db: db,
	}
}
func (r *favoriteRepository) CountFav(pctx context.Context, projectId primitive.ObjectID, userId string) (int64, int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	// $facet stage
	facetStage := bson.D{
		{"$facet", bson.M{
			"countUser": bson.A{
				bson.D{
					{"$match", bson.M{
						"$and": bson.A{
							bson.M{"_id": projectId},
							bson.M{"users": bson.M{"$elemMatch": bson.M{"$eq": userId}}},
						},
					}},
				},
				bson.D{
					{"$count", "countProject"},
				},
			},
			"countProject": bson.A{
				bson.D{
					{"$match", bson.M{
						"_id": projectId,
					}},
				},
				bson.D{
					{"$count", "countProject"},
				},
			},
		}},
	}

	// $project stage
	projectStage := bson.D{
		{"$project", bson.M{
			"result": bson.M{
				"$mergeObjects": bson.A{
					bson.M{"countUser": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$countUser.countProject", 0}}, 0}}},
					bson.M{"countProject": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$countProject.countProject", 0}}, 0}}},
				},
			},
		}},
	}

	// $replaceRoot stage
	replaceRootStage := bson.D{
		{"$replaceRoot", bson.M{"newRoot": "$result"}},
	}

	pipeline := mongo.Pipeline{facetStage, projectStage, replaceRootStage}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, errors.New("error: aggregation failed")
	}
	defer cursor.Close(ctx)

	var result struct {
		CountUser    int64 `bson:"countUser"`
		CountProject int64 `bson:"countProject"`
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return 0, 0, errors.New("error: decoding result failed")
		}
	}

	return result.CountProject, result.CountUser, nil
}
func (r *favoriteRepository) PushUserToFav(pctx context.Context, projectId primitive.ObjectID, userId string) (string, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	filter := bson.M{"_id": projectId}
	update := bson.M{"$push": bson.M{"users": userId}, "$set": bson.M{"updated_at": utils.LocalTime()}}

	result, err := col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error: Push User To Fav Failed: %s", err.Error())
		return "push", errors.New("error: push user to fav fail")
	}

	if result.ModifiedCount == 0 {
		return "push", errors.New("error: failed to push userid to fav")
	}

	return "push", nil
}

func (r *favoriteRepository) InsertOneFav(pctx context.Context, req *favorite.FavModel) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	_, err := col.InsertOne(context.Background(), req)
	if err != nil {
		log.Printf("Error: Insert One Fav Failed: %s", err.Error())
		return errors.New("error: insert one fav fail")
	}

	return nil
}

func (r *favoriteRepository) PullUserToFav(pctx context.Context, projectId primitive.ObjectID, userId string) (string, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	update := bson.M{"$pull": bson.M{"users": userId}}
	result, err := col.UpdateOne(context.Background(), bson.M{"_id": projectId}, update)
	if err != nil {
		log.Printf("Error: Pull User To Fav Failed: %s", err.Error())
		return "pull", errors.New("error: pull user to fav fail")
	}

	if result.ModifiedCount == 0 {
		return "pull", errors.New("error: failed to pull userid to fav")
	}

	return "pull", nil
}

// func ()
