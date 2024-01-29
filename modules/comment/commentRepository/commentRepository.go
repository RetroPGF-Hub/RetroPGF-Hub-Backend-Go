package commentrepository

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
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
	CommentRepositoryService interface {
		InsertEmptyComment(pctx context.Context, req *comment.CommentModel) error
		PushComment(pctx context.Context, projectId primitive.ObjectID, req *comment.CommentA) error
		PullComment(pctx context.Context, projectId primitive.ObjectID, commentId primitive.ObjectID) error
		CountComment(pctx context.Context, projectId primitive.ObjectID) (int64, error)
		UpdateComment(pctx context.Context, projectId primitive.ObjectID, req *comment.CommentA) (*comment.CommentModel, error)
	}

	commentRepository struct {
		db *mongo.Client
	}
)

func (r *commentRepository) commentDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("comment_db")
}

func NewCommentRepository(db *mongo.Client) CommentRepositoryService {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) InsertEmptyComment(pctx context.Context, req *comment.CommentModel) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.commentDbConn(ctx)
	col := db.Collection("comments")

	_, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: Insert One Comment Failed: %s", err.Error())
		return errors.New("error: insert one comment fail")
	}
	return nil

}

func (r *commentRepository) PushComment(pctx context.Context, projectId primitive.ObjectID, req *comment.CommentA) error {
	ctx, cancel := context.WithTimeout(pctx, 15*time.Second)
	defer cancel()

	db := r.commentDbConn(ctx)
	col := db.Collection("comments")

	filter := bson.M{"_id": projectId}
	update := bson.M{"$push": bson.M{"comments": req}, "$set": bson.M{"updated_at": utils.LocalTime()}}
	result, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error: Push Comment Failed: %s", err.Error())
		return errors.New("error: push comment fail")
	}

	if result.ModifiedCount == 0 {
		return errors.New("error: failed to push comment no document found")
	}
	return nil
}

func (r *commentRepository) PullComment(pctx context.Context, projectId primitive.ObjectID, commentId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(pctx, 15*time.Second)
	defer cancel()
	db := r.commentDbConn(ctx)
	col := db.Collection("comments")

	filter := bson.M{"_id": projectId}
	update := bson.M{
		"$pull": bson.M{"comments": bson.M{"_id": commentId}},
		"$set":  bson.M{"updated_at": utils.LocalTime()},
	}

	result, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error: Pull Comment Failed: %s", err.Error())
		return errors.New("error: pull comment fail")
	}

	if result.ModifiedCount == 0 {
		return errors.New("error: comment not found or failed to pull comment")
	}
	return nil
}

func (r *commentRepository) UpdateComment(pctx context.Context, projectId primitive.ObjectID, req *comment.CommentA) (*comment.CommentModel, error) {
	ctx, cancel := context.WithTimeout(pctx, 15*time.Second)
	defer cancel()

	db := r.commentDbConn(ctx)
	col := db.Collection("comments")

	filter := bson.M{"_id": projectId, "comments._id": req.CommentId, "comments.created_by": req.CreatedBy}
	update := bson.M{
		"$set": bson.M{
			"updated_at":            utils.LocalTime(),
			"comments.$.title":      req.Title,
			"comments.$.content":    req.Content,
			"comments.$.updated_at": utils.LocalTime(),
		},
	}
	updatedComment := new(comment.CommentModel)
	err := col.FindOneAndUpdate(ctx, filter, update).Decode(&updatedComment)
	if err != nil {
		log.Printf("Error: Update Comment Failed: %s", err.Error())
		return nil, errors.New("error: update comment failed")
	}

	return updatedComment, nil
}

// not sure if this query work fine or not
func (r *commentRepository) CountComment(pctx context.Context, projectId primitive.ObjectID) (int64, error) {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.commentDbConn(ctx)
	col := db.Collection("comments")

	filter := bson.M{"_id": projectId}
	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return count, errors.New("error: checking user unique failed")
	}
	return count, nil
}
