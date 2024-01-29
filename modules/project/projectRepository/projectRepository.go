package projectrepository

import (
	commentPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentPb"
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favPb"
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
		InsertOneProject(pctx context.Context, req *project.ProjectModel, favGrpcUrl string) (primitive.ObjectID, error)
		FindOneProject(pctx context.Context, projectId primitive.ObjectID) (*project.ProjectModel, error)
		FindOneUserWithId(pctx context.Context, grpcUrl string, req *usersPb.GetUserInfoReq) (*usersPb.GetUserInfoRes, error)
		UpdateProject(pctx context.Context, req *project.ProjectModel, userId string) (*project.ProjectModel, error)
		DeleteProject(pctx context.Context, projectId primitive.ObjectID, userId string) error
		InsertEmptyFav(pctx context.Context, projectId primitive.ObjectID, favGrpcUrl string) error
		InsertEmptyComment(pctx context.Context, projectId primitive.ObjectID, commentGrpcUrl string) error
		DeleteFavProject(pctx context.Context, projectId primitive.ObjectID, grpcUrl string) (*favPb.DeleteFavProjectRes, error)
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

func (r *projectRepository) InsertOneProject(pctx context.Context, req *project.ProjectModel, favGrpcUrl string) (primitive.ObjectID, error) {

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

func (r *projectRepository) FindOneProject(pctx context.Context, projectId primitive.ObjectID) (*project.ProjectModel, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("projects")

	result := new(project.ProjectModel)

	if err := col.FindOne(ctx, bson.M{"_id": projectId}).Decode(result); err != nil {
		return nil, errors.New("error: project id not found")
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

// make grpc conn to create fav document # project -> fav
func (r *projectRepository) InsertEmptyFav(pctx context.Context, projectId primitive.ObjectID, favGrpcUrl string) error {
	jwtauth.SetApiKeyInContext(&pctx)
	conn, err := grpcconn.NewGrpcClient(favGrpcUrl)
	if err != nil {
		log.Printf("Error: Grpc Conn Error %s", err.Error())
		return errors.New("error: grpc connection to fav failed")
	}

	result, err := conn.Fav().CreateFavProject(pctx, &favPb.CreateFavProjectReq{
		ProjectId: projectId.Hex(),
	})
	if err != nil {
		log.Printf("Error: Create Fav For Project Error %s", err.Error())
		return errors.New(result.Msg)
	}

	return nil
}

// make grpc conn to create fav document # project -> fav
func (r *projectRepository) InsertEmptyComment(pctx context.Context, projectId primitive.ObjectID, commentGrpcUrl string) error {
	jwtauth.SetApiKeyInContext(&pctx)
	conn, err := grpcconn.NewGrpcClient(commentGrpcUrl)
	if err != nil {
		log.Printf("Error: Grpc Conn Error %s", err.Error())
		return errors.New("error: grpc connection to comment failed")
	}

	_, err = conn.Comment().CreateCommentProject(pctx, &commentPb.CreateCommentProjectReq{
		ProjectId: projectId.Hex(),
	})

	if err != nil {
		log.Printf("Error: Create Comment For Project Error %s", err.Error())
		return err
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
		return nil, errors.New("error: find out user for project not found")
	}

	return result, nil
}

// make grpc conn to delete fav # project -> fav
// if some of the process cancel we need to back track and delete all of the thing
func (r *projectRepository) DeleteFavProject(pctx context.Context, projectId primitive.ObjectID, grpcUrl string) (*favPb.DeleteFavProjectRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpcconn.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: Grpc Conn Error %s", err.Error())
		return nil, errors.New("error: grpc connection failed")
	}
	result, err := conn.Fav().DeleteFavProject(ctx, &favPb.DeleteFavProjectReq{
		ProjectId: projectId.Hex(),
	})

	if err != nil {
		log.Printf("Error: Delete Fav Project Grpc Failed %s", err.Error())
		return nil, errors.New("error: delete fav project grpc failed")
	}

	return result, nil
}
