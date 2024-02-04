package migration

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/database"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func commentDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("comment_db")
}

func Commentigrate(pctx context.Context, cfg *config.Config) {
	db := commentDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	// comment
	col := db.Collection("comments")

	// indexs comment
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"_id", 1}},
		},
		{
			Keys: bson.D{{"_id", 1}, {"comments._id", 1}, {"comments.created_by", 1}},
		},
	})

	for _, index := range indexs {
		log.Printf("Indexs: %s", index)
	}
}
