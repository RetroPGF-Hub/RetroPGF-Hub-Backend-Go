package migration

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/database"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func usersDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("user_db")
}

func UsersMigrate(pctx context.Context, cfg *config.Config) {
	db := usersDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	// users
	col := db.Collection("users")

	// indexs users
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"_id", 1}},
		},
		{
			Keys: bson.D{{"email", 1}},
		},
	})

	for _, index := range indexs {
		log.Printf("Indexs: %s", index)
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

	document := func() []any {
		users := []*users.UserDb{
			{
				Id:        primitive.NewObjectID(),
				Source:    "web",
				Email:     "mix123@gmail.com",
				Profile:   "https://google.com",
				Password:  string(hashPassword),
				Username:  "mix",
				Firstname: "mixfirstname",
				Lastname:  "mixlastname",
				CreateAt:  utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
		}

		docs := make([]any, 0)
		for _, r := range users {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(pctx, document, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate users complete: ", results)

}
