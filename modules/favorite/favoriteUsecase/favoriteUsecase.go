package favoriteusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoritePb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	FavoriteUsecaseService interface {
		FavPullOrPushUsecase(pctx context.Context, projectId, userId string) (string, error)
		GetAllProjectByUserId(pctx context.Context, req *favPb.GetAllFavReq) (*favPb.GetAllFavRes, error)
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

func (u *favoriteUsecase) GetAllProjectByUserId(pctx context.Context, req *favPb.GetAllFavReq) (*favPb.GetAllFavRes, error) {

	data, err := u.pActor.FavoriteRepo.GetAllProjectInUser(pctx, utils.ConvertToObjectId(req.UserId))
	if err != nil {
		return nil, err
	}

	var projectIds []primitive.ObjectID

	for _, v := range data.ProjectId {
		projectIds = append(projectIds, utils.ConvertToObjectId(v))
	}
	projects, err := u.pActor.ProjectRepo.FindManyProjectId(pctx, projectIds)
	if err != nil {
		return nil, err
	}

	projectsRes := make([]*favPb.ProjectRes, 0)

	for _, v := range projects {
		projectsRes = append(projectsRes, &favPb.ProjectRes{
			Id:             v.Id.Hex(),
			Name:           v.Name,
			BannerUrl:      v.BannerUrl,
			WebsiteUrl:     v.WebsiteUrl,
			CryptoCategory: v.CryptoCategory,
			Description:    v.Description,
			Reason:         v.Reason,
			Category:       v.Category,
			Contact:        v.Contact,
			FavCount:       v.FavCount,
			CommentCount:   v.CommentCount,
			CreatedBy:      v.CreatedBy,
			CreatedAt:      v.CreateAt.String(),
		})
	}

	return &favPb.GetAllFavRes{
		UserId:   data.User.Hex(),
		Projects: projectsRes,
	}, nil

}
