package datacenterrepository

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
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
	DatacenterRepositoryService interface {
		GetAllProjectRepo(pctx context.Context, limit, skip int64) ([]*project.ProjectModel, error)
		GetSingleProjectRepo(pctx context.Context, projectId primitive.ObjectID) (*project.ProjectModel, error)
		InsertUrlCache(pctx context.Context, req *datacenter.CacheModel) (primitive.ObjectID, error)
		DeleteUrlCache(pctx context.Context, urlId primitive.ObjectID) error
		GetAllUrlCache(pctx context.Context) ([]*datacenter.CacheModel, error)
	}

	datacenterRepository struct {
		db *mongo.Client
	}
)

func NewDatacenterRepository(db *mongo.Client) DatacenterRepositoryService {
	return &datacenterRepository{
		db: db,
	}
}

func (r *datacenterRepository) projectDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("project_db")
}

func (r *datacenterRepository) cacheDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("cache_db")
}

func (r *datacenterRepository) GetAllProjectRepo(pctx context.Context, limit, skip int64) ([]*project.ProjectModel, error) {

	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()
	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	// descending order
	findOptions := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{"created_at", -1}})

	cursor, err := col.Find(pctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(pctx)

	var projects []*project.ProjectModel
	if err := cursor.All(pctx, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *datacenterRepository) GetSingleProjectRepo(pctx context.Context, projectId primitive.ObjectID) (*project.ProjectModel, error) {

	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()
	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	projects := new(project.ProjectModel)

	if err := col.FindOne(pctx, bson.M{"_id": projectId}).Decode(&projects); err != nil {
		return nil, errors.New("error: get single project datacenter failed")
	}

	return projects, nil
}

func (r *datacenterRepository) InsertUrlCache(pctx context.Context, req *datacenter.CacheModel) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()
	db := r.cacheDbConn(ctx)
	col := db.Collection("cache_db")

	cacheId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertUrlCache: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert url repo fail")
	}

	return cacheId.InsertedID.(primitive.ObjectID), nil
}

func (r *datacenterRepository) GetAllUrlCache(pctx context.Context) ([]*datacenter.CacheModel, error) {

	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()
	db := r.cacheDbConn(ctx)
	col := db.Collection("cache_db")

	findOptions := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := col.Find(pctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(pctx)

	var urls []*datacenter.CacheModel
	if err := cursor.All(pctx, &urls); err != nil {
		return nil, err
	}

	return urls, nil

}

func (r *datacenterRepository) DeleteUrlCache(pctx context.Context, urlId primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()
	db := r.cacheDbConn(ctx)
	col := db.Collection("cache_db")

	result, err := col.DeleteOne(ctx, bson.M{"_id": urlId})
	if err != nil {
		log.Printf("Error: Delete Url Repo: %s", err.Error())
		return errors.New("error: delete url repo fail")
	}

	if result.DeletedCount == 0 {
		return errors.New("error: no dument found to delete or access denied")
	}

	return nil
}
