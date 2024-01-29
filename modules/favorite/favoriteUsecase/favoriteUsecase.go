package favoriteusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite"
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favPb"
	favoriterepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"
)

type (
	FavoriteUsecaseService interface {
		FavPullOrPushUsecase(pctx context.Context, projectId, userId string) (string, error)
		InsertEmptyDoc(pctx context.Context, req *favPb.CreateFavProjectReq) error
		DeleteFavUsecase(pctx context.Context, req *favPb.DeleteFavProjectReq) error
	}

	favoriteUsecase struct {
		favoriteRepo favoriterepository.FavoriteRepositoryService
	}
)

func NewFavoriteUsecase(favoriteRepo favoriterepository.FavoriteRepositoryService) FavoriteUsecaseService {
	return &favoriteUsecase{
		favoriteRepo: favoriteRepo,
	}
}

func (u *favoriteUsecase) FavPullOrPushUsecase(pctx context.Context, projectId, userId string) (string, error) {
	projectIdPri := utils.ConvertToObjectId(projectId)

	countP, countU, err := u.favoriteRepo.CountFav(pctx, projectIdPri, userId)
	if err != nil {
		return "", err
	}

	if countP == 0 {
		// if err := u.favoriteRepo.InsertOneFav(pctx, &favorite.FavModel{
		// 	ProjectId: projectIdPri,
		// 	Users:     make([]string, 0),
		// 	CreateAt:  utils.LocalTime(),
		// 	UpdatedAt: utils.LocalTime(),
		// }); err != nil {
		// 	return "", err
		// }
		return "", errors.New("project you looking for is not exist")
	}
	var opera string

	if countU == 0 {
		push, err := u.favoriteRepo.PushUserToFav(pctx, projectIdPri, userId)
		if err != nil {
			return push, err
		}
		opera = push
	} else {
		push, err := u.favoriteRepo.PullUserToFav(pctx, projectIdPri, userId)
		if err != nil {
			return push, err
		}
		opera = push
	}
	return opera, nil
}

// For grpc from project
func (u *favoriteUsecase) InsertEmptyDoc(pctx context.Context, req *favPb.CreateFavProjectReq) error {
	if err := u.favoriteRepo.InsertOneFav(pctx, &favorite.FavModel{
		ProjectId: utils.ConvertToObjectId(req.ProjectId),
		Users:     make([]string, 0),
		CreateAt:  utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
	}); err != nil {
		return err
	}
	return nil
}

// for grpc cancel project
func (u *favoriteUsecase) DeleteFavUsecase(pctx context.Context, req *favPb.DeleteFavProjectReq) error {
	if err := u.favoriteRepo.DeleteFav(pctx, utils.ConvertToObjectId(req.ProjectId)); err != nil {
		return err
	}
	return nil
}
