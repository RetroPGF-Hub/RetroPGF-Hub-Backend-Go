package projectrepository

import (
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

func (r *projectRepository) InsertOneQuestion(pctx context.Context, req *project.QuestionModel) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()
	db := r.projectDbConn(ctx)
	col := db.Collection("questions")

	questionId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneProject: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one project fail")
	}
	return questionId.InsertedID.(primitive.ObjectID), nil

}

func (r *projectRepository) FindOneQuestion(pctx context.Context, questionId string) (*project.QuestionModel, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("questions")

	questions := new(project.QuestionModel)

	if err := col.FindOne(pctx, bson.M{"_id": questionId}).Decode(&questions); err != nil {
		return nil, errors.New("error: get single question failed")
	}

	return questions, nil
}

func (r *projectRepository) UpdateQuestionFavCount(pctx context.Context, questionId primitive.ObjectID, counter int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("questions")

	filter := bson.M{"_id": questionId}

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
		log.Printf("Error: Update Fav Count Question Error %s", err.Error())
		return errors.New("error: update fav count question failed")
	}

	if result.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (r *projectRepository) UpdateQuestionCommentCount(pctx context.Context, questionId primitive.ObjectID, counter int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.projectDbConn(ctx)
	col := db.Collection("questions")

	filter := bson.M{"_id": questionId}

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
		log.Printf("Error: Update Comment Count Question Error %s", err.Error())
		return errors.New("error: update comment count failed")
	}

	if result.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (r *projectRepository) FindAllQuestions(pctx context.Context, grpcUrl string, limit, skip int64) ([]*project.QuestionModel, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	// jwtauth.SetApiKeyInContext(&ctx)
	// conn, err := grpcconn.NewGrpcClient(grpcUrl)
	// if err != nil {
	// 	log.Printf("Error: Grpc Conn Error %s", err.Error())
	// 	return nil, errors.New("error: grpc connection failed")
	// }
	// result, err := conn.Datacenter().GetProjectDataCenter(ctx, req)
	// if err != nil {
	// 	log.Printf("Error: Find All Project Error %s", err.Error())
	// 	return nil, errors.New("error: find all project failed")
	// }

	db := r.projectDbConn(ctx)
	col := db.Collection("questions")

	// descending order
	findOptions := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{"created_at", -1}})

	cursor, err := col.Find(pctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(pctx)

	var questions []*project.QuestionModel
	if err := cursor.All(pctx, &questions); err != nil {
		return nil, err
	}

	return questions, nil

}
