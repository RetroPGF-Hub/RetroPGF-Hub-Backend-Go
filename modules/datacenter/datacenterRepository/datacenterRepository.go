package datacenterrepository

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"context"
	"errors"
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
