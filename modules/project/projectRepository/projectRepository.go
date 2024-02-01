package projectrepository

import (
	datacenterPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterPb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
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
	ProjectRepositoryService interface {
		InsertOneProject(pctx context.Context, req *project.ProjectModel) (primitive.ObjectID, error)
		FindOneProject(pctx context.Context, projectId, datacenterUrl string) (*datacenterPb.GetSingleProjectDataCenterRes, error)
		FindOneUserWithId(pctx context.Context, grpcUrl string, req *usersPb.GetUserInfoReq) (*usersPb.GetUserInfoRes, error)
		UpdateProject(pctx context.Context, req *project.ProjectModel, userId string) (*project.ProjectModel, error)
		DeleteProject(pctx context.Context, projectId primitive.ObjectID, userId string) error
		UpdateFavCount(pctx context.Context, projectId primitive.ObjectID, counter int64) error
		UpdateCommentCount(pctx context.Context, projectId primitive.ObjectID, counter int64) error
		FindManyProjectId(pctx context.Context, projectId []primitive.ObjectID) ([]*project.ProjectModel, error)
		FindAllProjectDatacenter(pctx context.Context, grpcUrl string, req *datacenterPb.GetProjectDataCenterReq) (*datacenterPb.GetProjectDataCenterRes, error)
	}

	projectRepository struct {
		db *mongo.Client
	}
)

func (r *projectRepository) projectDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("project_db")
}

func NewProjectRepository(db *mongo.Client) ProjectRepositoryService {
	return &projectRepository{
		db: db,
	}
}

func (r *projectRepository) InsertOneProject(pctx context.Context, req *project.ProjectModel) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()
	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	projectId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneProject: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one project fail")
	}
	return projectId.InsertedID.(primitive.ObjectID), nil

}

func (r *projectRepository) FindOneProject(pctx context.Context, projectId, datacenterUrl string) (*datacenterPb.GetSingleProjectDataCenterRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpcconn.NewGrpcClient(datacenterUrl)
	if err != nil {
		log.Printf("Error: Grpc Conn Error %s", err.Error())
		return nil, errors.New("error: grpc connection failed")
	}
	result, err := conn.Datacenter().GetSingleProjectDataCenter(ctx, &datacenterPb.GetSingleProjectDataCenterReq{ProjecId: projectId})
	if err != nil {
		log.Printf("Error: Find One Project Error %s", err.Error())
		return nil, errors.New("error: find one project failed")
	}

	return result, nil
}

func (r *projectRepository) UpdateProject(pctx context.Context, req *project.ProjectModel, userId string) (*project.ProjectModel, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	filter := bson.M{"_id": req.Id, "created_by": userId}
	update := bson.M{
		"$set": bson.M{
			"name":            req.Name,
			"description":     req.Description,
			"logo_url":        req.LogoUrl,
			"banner_url":      req.BannerUrl,
			"website_url":     req.WebsiteUrl,
			"crypto_category": req.CryptoCategory,
			"reason":          req.Reason,
			"category":        req.Category,
			"contact":         req.Contact,
			"updated_at":      time.Now(),
		},
	}

	options := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)

	updatedProject := new(project.ProjectModel)
	err := col.FindOneAndUpdate(ctx, filter, update, options).Decode(updatedProject)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("warning: no document found")
		}
		log.Printf("Error: FindOneAndUpdate Project Error %s", err.Error())
		return nil, errors.New("error: find one and update project failed or access denied")
	}

	return updatedProject, nil
}

func (r *projectRepository) UpdateFavCount(pctx context.Context, projectId primitive.ObjectID, counter int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	filter := bson.M{"_id": projectId}

	update := bson.M{
		"$inc": bson.M{
			"fav_count": counter,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	options := options.UpdateOptions(*options.Update().SetUpsert(false))

	result, err := col.UpdateOne(ctx, filter, update, &options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("warning: no document found")
		}
		log.Printf("Error: Update Fav Count Project Error %s", err.Error())
		return errors.New("error: update fav count failed")
	}

	if result.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (r *projectRepository) UpdateCommentCount(pctx context.Context, projectId primitive.ObjectID, counter int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	filter := bson.M{"_id": projectId}

	update := bson.M{
		"$inc": bson.M{
			"comment_count": counter,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	options := options.UpdateOptions(*options.Update().SetUpsert(false))

	result, err := col.UpdateOne(ctx, filter, update, &options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("warning: no document found")
		}
		log.Printf("Error: Update Comment Count Project Error %s", err.Error())
		return errors.New("error: update comment count failed")
	}

	if result.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (r *projectRepository) DeleteProject(pctx context.Context, projectId primitive.ObjectID, userId string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	result, err := col.DeleteOne(ctx, bson.M{"_id": projectId, "created_by": userId})
	if err != nil {
		log.Printf("Error: Delete One Project Error %s", err.Error())
		return errors.New("error: delete one project failed")
	}

	if result.DeletedCount == 0 {
		return errors.New("error: no dument found to delete or access denied")
	}

	return nil

}

// make grpc conn to get info of user # project -> user
func (r *projectRepository) FindOneUserWithId(pctx context.Context, grpcUrl string, req *usersPb.GetUserInfoReq) (*usersPb.GetUserInfoRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpcconn.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: Grpc Conn Error %s", err.Error())
		return nil, errors.New("error: grpc connection failed")
	}
	result, err := conn.Users().GetUserInfoById(ctx, req)
	if err != nil {
		log.Printf("Error: Find One User For Project Error %s", err.Error())
		return nil, errors.New("error: find one user for project not found")
	}

	return result, nil
}

func (r *projectRepository) FindManyProjectId(pctx context.Context, projectId []primitive.ObjectID) ([]*project.ProjectModel, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	matchInState := bson.D{
		{"$match",
			bson.D{
				{"_id",
					bson.D{
						{"$in",
							projectId,
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
				{"name", 1},
				{"banner_url", 1},
				{"logo_url", 1},
				{"website_url", 1},
				{"crypto_category", 1},
				{"description", 1},
				{"reason", 1},
				{"category", 1},
				{"contact", 1},
				{"fav_count", 1},
				{"comment_count", 1},
				{"created_by", 1},
				{"created_at", 1},
			},
		},
	}

	pipeline := mongo.Pipeline{matchInState, projectInState}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*project.ProjectModel

	for cursor.Next(ctx) {
		var result project.ProjectModel
		err := cursor.Decode(&result)
		if err != nil {
			return nil, errors.New("error: decoding result failed")
		}
		projects = append(projects, &result)
	}

	return projects, nil

}

func (r *projectRepository) FindAllProjectDatacenter(pctx context.Context, grpcUrl string, req *datacenterPb.GetProjectDataCenterReq) (*datacenterPb.GetProjectDataCenterRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpcconn.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: Grpc Conn Error %s", err.Error())
		return nil, errors.New("error: grpc connection failed")
	}
	result, err := conn.Datacenter().GetProjectDataCenter(ctx, req)
	if err != nil {
		log.Printf("Error: Find All Project Error %s", err.Error())
		return nil, errors.New("error: find all project failed")
	}

	return result, nil
}
