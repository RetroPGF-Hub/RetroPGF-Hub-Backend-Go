package datacenterusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	datacenterPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterPb"
	datacenterrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
)

type (
	DatacenterUsecaseService interface {
		GetAllProjectUsecase(pctx context.Context, limit, skip int64) (*datacenterPb.GetProjectDataCenterRes, error)
		GetSingleProjectUsecase(pctx context.Context, projectId string) (*datacenterPb.GetSingleProjectDataCenterRes, error)
	}

	datacenterUsecase struct {
		cfg            *config.Grpc
		datacenterRepo datacenterrepository.DatacenterRepositoryService
	}
)

func NewDatacenterUsecase(datacenterRepo datacenterrepository.DatacenterRepositoryService, cfg *config.Grpc) DatacenterUsecaseService {
	return &datacenterUsecase{
		datacenterRepo: datacenterRepo,
		cfg:            cfg,
	}
}

func (u *datacenterUsecase) GetAllProjectUsecase(pctx context.Context, limit, skip int64) (*datacenterPb.GetProjectDataCenterRes, error) {
	projects, err := u.datacenterRepo.GetAllProjectRepo(pctx, limit, skip)
	if err != nil {
		return nil, err
	}

	var result []*datacenterPb.ProjectRes

	for _, v := range projects {
		result = append(result, &datacenterPb.ProjectRes{
			Id:             v.Id.Hex(),
			Name:           v.Name,
			LogoUrl:        v.LogoUrl,
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
			UpdatedAt:      v.UpdatedAt.String(),
		})
	}

	return &datacenterPb.GetProjectDataCenterRes{
		Projects: result,
	}, nil

}

func (u *datacenterUsecase) GetSingleProjectUsecase(pctx context.Context, projectId string) (*datacenterPb.GetSingleProjectDataCenterRes, error) {
	projects, err := u.datacenterRepo.GetSingleProjectRepo(pctx, utils.ConvertToObjectId(projectId))
	if err != nil {
		return nil, err
	}

	return &datacenterPb.GetSingleProjectDataCenterRes{
		Projects: &datacenterPb.ProjectRes{
			Id:             projects.Id.Hex(),
			Name:           projects.Name,
			LogoUrl:        projects.LogoUrl,
			BannerUrl:      projects.BannerUrl,
			WebsiteUrl:     projects.WebsiteUrl,
			CryptoCategory: projects.CryptoCategory,
			Description:    projects.Description,
			Reason:         projects.Reason,
			Category:       projects.Category,
			Contact:        projects.Contact,
			FavCount:       projects.FavCount,
			CommentCount:   projects.CommentCount,
			CreatedBy:      projects.CreatedBy,
			CreatedAt:      projects.CreateAt.String(),
			UpdatedAt:      projects.UpdatedAt.String(),
		},
	}, nil

}
