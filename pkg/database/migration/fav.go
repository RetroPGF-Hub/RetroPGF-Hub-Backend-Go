package migration

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/database"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func favDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("favorite_db")
}

func FavMigrate(pctx context.Context, cfg *config.Config) {
	db := favDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	// fav
	col := db.Collection("favs")

	// indexs fav
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"_id", 1}},
		},
	})

	for _, index := range indexs {
		log.Printf("Indexs: %s", index)
	}
}
