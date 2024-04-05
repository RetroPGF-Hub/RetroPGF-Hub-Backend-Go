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
		PushProjectToFav(pctx context.Context, projectId string, userId primitive.ObjectID) (string, error)
		PullProjectToFav(pctx context.Context, projectId string, userId primitive.ObjectID) (string, error)
		CountFav(pctx context.Context, userId primitive.ObjectID, projectId string) (int64, int64, error)
		InsertOneFav(pctx context.Context, req *favorite.FavProjectModel) error
		DeleteFav(pctx context.Context, userId primitive.ObjectID) error
		CountUserFav(pctx context.Context, userId primitive.ObjectID) (int64, error)
		GetAllProjectInUser(pctx context.Context, userId primitive.ObjectID) (*favorite.FavProjectModel, error)
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
func (r *favoriteRepository) CountFav(pctx context.Context, userId primitive.ObjectID, projectId string) (int64, int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	// $facet stage
	facetStage := bson.D{
		{"$facet", bson.M{
			"countProject": bson.A{
				bson.D{
					{"$match", bson.M{
						"$and": bson.A{
							bson.M{"_id": userId},
							bson.M{"projects": bson.M{"$elemMatch": bson.M{"$eq": projectId}}},
						},
					}},
				},
				bson.D{
					{"$count", "countProject"},
				},
			},
			"countUser": bson.A{
				bson.D{
					{"$match", bson.M{
						"_id": userId,
					}},
				},
				bson.D{
					{"$count", "countUser"},
				},
			},
		}},
	}

	// $project stage
	projectStage := bson.D{
		{"$project", bson.M{
			"result": bson.M{
				"$mergeObjects": bson.A{
					bson.M{"countProject": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$countProject.countProject", 0}}, 0}}},
					bson.M{"countUser": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$countUser.countUser", 0}}, 0}}},
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

func (r *favoriteRepository) CountUserFav(pctx context.Context, userId primitive.ObjectID) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	count, err := col.CountDocuments(pctx, bson.M{"_id": userId})
	if err != nil {
		return -1, errors.New("error: count user fav failed")
	}

	return count, nil
}

func (r *favoriteRepository) PushProjectToFav(pctx context.Context, projectId string, userId primitive.ObjectID) (string, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	filter := bson.M{"_id": userId}
	update := bson.M{"$push": bson.M{"projects": projectId}, "$set": bson.M{"updated_at": utils.LocalTime()}}

	result, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error: Push User To Fav Failed: %s", err.Error())
		return "push", errors.New("error: push user to fav fail")
	}

	if result.ModifiedCount == 0 {
		return "push", errors.New("error: failed to push userid to fav")
	}

	return "push", nil
}

func (r *favoriteRepository) PullProjectToFav(pctx context.Context, projectId string, userId primitive.ObjectID) (string, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	update := bson.M{"$pull": bson.M{"projects": projectId}}
	result, err := col.UpdateOne(ctx, bson.M{"_id": userId}, update)
	if err != nil {
		log.Printf("Error: Pull User To Fav Failed: %s", err.Error())
		return "pull", errors.New("error: pull user to fav fail")
	}

	if result.ModifiedCount == 0 {
		return "pull", errors.New("error: failed to pull userid to fav")
	}

	return "pull", nil
}

func (r *favoriteRepository) InsertOneFav(pctx context.Context, req *favorite.FavProjectModel) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	_, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: Insert One Fav Failed: %s", err.Error())
		return errors.New("error: insert one fav fail")
	}

	return nil
}

func (r *favoriteRepository) DeleteFav(pctx context.Context, userId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")

	result, err := col.DeleteOne(ctx, bson.M{"_id": userId})
	if err != nil {
		log.Printf("Error: Delete One Fav Error %s", err.Error())
		return errors.New("error: delete one fav failed")
	}

	if result.DeletedCount == 0 {
		return errors.New("error: no dument found to delete")
	}

	return nil
}

func (r *favoriteRepository) GetAllProjectInUser(pctx context.Context, userId primitive.ObjectID) (*favorite.FavProjectModel, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.favoriteDbConn(ctx)
	col := db.Collection("favs")
	filter := bson.M{"_id": userId}

	result := new(favorite.FavProjectModel)

	if err := col.FindOne(ctx, filter).Decode(result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &favorite.FavProjectModel{}, nil
		} else {
			return nil, errors.New("error: get all fav id not found")
		}
	}

	return result, nil
}
