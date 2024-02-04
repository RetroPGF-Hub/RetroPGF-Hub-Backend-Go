package migration

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/database"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func projectDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("project_db")
}

func ProjectMigrate(pctx context.Context, cfg *config.Config) {
	db := projectDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	// project
	col := db.Collection("projects")

	// indexs project
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"_id", 1}},
		},
		{
			Keys: bson.D{{"_id", 1}, {"created_by", 1}},
		},
	})

	for _, index := range indexs {
		log.Printf("Indexs: %s", index)
	}

	document := func() []any {
		var projects []*project.ProjectModel

		for i := 0; i < 1; i++ {
			project := &project.ProjectModel{
				Id:             primitive.NewObjectID(),
				Name:           fmt.Sprintf("Project%d", i),
				LogoUrl:        fmt.Sprintf("https://example.com/logo%d.png", i),
				BannerUrl:      fmt.Sprintf("https://example.com/banner%d.png", i),
				WebsiteUrl:     fmt.Sprintf("https://project%d.com", i),
				CryptoCategory: "Blockchain",
				Description:    fmt.Sprintf("Description for Project%d", i),
				Reason:         fmt.Sprintf("Reason for Project%d", i),
				Category:       "Technology",
				Contact:        fmt.Sprintf("info@project%d.com", i),
				FavCount:       0,
				CommentCount:   0,
				CreatedBy:      "65b94263192fdc9412761367",
				CreateAt:       utils.LocalTime(),
				UpdatedAt:      utils.LocalTime(),
			}
			projects = append(projects, project)
		}

		docs := make([]any, 0)
		for _, r := range projects {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(pctx, document, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate project complete: ", results)

}
