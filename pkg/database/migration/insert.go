package migration

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func InsertTestData(pctx context.Context, cfg *config.Config) {
	dbProject := projectDbConn(pctx, cfg)
	defer dbProject.Client().Disconnect(pctx)
	dbUser := usersDbConn(pctx, cfg)
	defer dbUser.Client().Disconnect(pctx)
	dbFav := favDbConn(pctx, cfg)
	defer dbFav.Client().Disconnect(pctx)
	dbCom := commentDbConn(pctx, cfg)
	defer dbCom.Client().Disconnect(pctx)

	colUser := dbUser.Collection("users")
	colProject := dbProject.Collection("projects")
	colFav := dbFav.Collection("favs")
	colCom := dbCom.Collection("comments")

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

	dProject, dUser, dFav, dCom := func() ([]any, []any, []any, []any) {
		var rProjects []*project.ProjectModel
		var rUsers []*users.UserDb
		var rFavs []*favorite.FavModel
		var rComs []*comment.CommentModel

		for i := 0; i < 600; i++ {
			userId := primitive.NewObjectID()
			projectId := primitive.NewObjectID()
			user := &users.UserDb{
				Id:        userId,
				Source:    "web",
				Email:     fmt.Sprintf("mix%d@gmail.com", i),
				Profile:   "https://google.com",
				Password:  string(hashPassword),
				Username:  fmt.Sprintf("mix%d", i),
				Firstname: fmt.Sprintf("mix Firstname%d", i),
				Lastname:  fmt.Sprintf("mix LastName%d", i),
				CreateAt:  utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			}

			project := &project.ProjectModel{
				Id:             projectId,
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
				CreatedBy:      userId.Hex(),
				CreateAt:       utils.LocalTime(),
				UpdatedAt:      utils.LocalTime(),
			}

			fav := &favorite.FavModel{
				User:      userId,
				ProjectId: make([]string, 0),
				CreateAt:  utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			}

			comment := &comment.CommentModel{
				ProjectId: projectId,
				Comments:  []comment.CommentA{},
				CreateAt:  utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			}
			rProjects = append(rProjects, project)
			rUsers = append(rUsers, user)
			rFavs = append(rFavs, fav)
			rComs = append(rComs, comment)
		}

		docsProject := make([]any, 0)
		for _, r := range rProjects {
			docsProject = append(docsProject, r)
		}
		docsUser := make([]any, 0)
		for _, r := range rUsers {
			docsUser = append(docsUser, r)
		}
		docsFavs := make([]any, 0)
		for _, r := range rFavs {
			docsFavs = append(docsFavs, r)
		}
		docsComs := make([]any, 0)
		for _, r := range rComs {
			docsComs = append(docsComs, r)
		}
		return docsProject, docsUser, docsFavs, docsComs
	}()

	results, err := colProject.InsertMany(pctx, dProject, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate project complete: ", results)

	resultsU, err := colUser.InsertMany(pctx, dUser, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate users complete: ", resultsU)

	resultsF, err := colFav.InsertMany(pctx, dFav, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate users complete: ", resultsF)

	resultsC, err := colCom.InsertMany(pctx, dCom, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate users complete: ", resultsC)

}
