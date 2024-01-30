package favoriteusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"
)

type (
	FavoriteUsecaseService interface {
		FavPullOrPushUsecase(pctx context.Context, projectId, userId string) (string, error)
	}

	favoriteUsecase struct {
		pActor modules.ProjectSvcInteractor
	}
)

func NewFavoriteUsecase(pActor modules.ProjectSvcInteractor) FavoriteUsecaseService {
	return &favoriteUsecase{
		pActor: pActor,
	}
}

func (u *favoriteUsecase) FavPullOrPushUsecase(pctx context.Context, projectId, userId string) (string, error) {
	userIdPri := utils.ConvertToObjectId(userId)
	projectIdPri := utils.ConvertToObjectId(projectId)
	countP, countU, err := u.pActor.FavoriteRepo.CountFav(pctx, userIdPri, projectId)
	if err != nil {
		return "", err
	}

	if countU == 0 {
		return "", errors.New("user you looking for is not exist")
	}
	var opera string
	if countP == 0 {
		push, err := u.pActor.FavoriteRepo.PushProjectToFav(pctx, projectId, userIdPri)
		if err != nil {
			return push, err
		}

		if err := u.pActor.ProjectRepo.UpdateFavCount(pctx, projectIdPri, 1); err != nil {
			return push, err
		}

		opera = push
	} else {
		pull, err := u.pActor.FavoriteRepo.PullProjectToFav(pctx, projectId, userIdPri)
		if err != nil {
			return pull, err
		}
		if err := u.pActor.ProjectRepo.UpdateFavCount(pctx, projectIdPri, -1); err != nil {
			return pull, err
		}
		opera = pull
	}
	return opera, nil
}
