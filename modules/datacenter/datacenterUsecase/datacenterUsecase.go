package datacenterusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter"
	datacenterPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterPb"
	datacenterrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
)

type (
	DatacenterUsecaseService interface {
		GetAllProjectUsecase(pctx context.Context, limit, skip int64) (*datacenterPb.GetProjectDataCenterRes, error)
		GetSingleProjectUsecase(pctx context.Context, projectId string) (*datacenterPb.GetSingleProjectDataCenterRes, error)
		InsertUrlCache(pctx context.Context, url string) (string, error)
		DeletetUrlCahce(pctx context.Context, url string) error
		FindManyUrlsCache(pctx context.Context) ([]*datacenter.CacheModel, error)
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
		result = append(result, u.convertToPbProject(v))
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
		Projects: u.convertToPbProject(projects),
	}, nil

}

func (u *datacenterUsecase) convertToPbProject(p *project.ProjectModel) *datacenterPb.ProjectRes {
	return &datacenterPb.ProjectRes{
		Id:             p.Id.Hex(),
		Name:           p.Name,
		LogoUrl:        p.LogoUrl,
		BannerUrl:      p.BannerUrl,
		WebsiteUrl:     p.WebsiteUrl,
		CryptoCategory: p.CryptoCategory,
		Description:    p.Description,
		Reason:         p.Reason,
		Category:       p.Category,
		Contact:        p.Contact,
		FavCount:       p.FavCount,
		CommentCount:   p.CommentCount,
		CreatedBy:      p.CreatedBy,
		CreatedAt:      p.CreateAt.String(),
		UpdatedAt:      p.UpdatedAt.String(),
	}
}

func (u *datacenterUsecase) InsertUrlCache(pctx context.Context, url string) (string, error) {
	id, err := u.datacenterRepo.InsertUrlCache(pctx, &datacenter.CacheModel{Url: url})
	if err != nil {
		return "", err
	}
	return id.Hex(), err
}

func (u *datacenterUsecase) DeletetUrlCahce(pctx context.Context, url string) error {
	err := u.datacenterRepo.DeleteUrlCache(pctx, utils.ConvertToObjectId(url))
	if err != nil {
		return err
	}
	return nil
}

func (u *datacenterUsecase) FindManyUrlsCache(pctx context.Context) ([]*datacenter.CacheModel, error) {
	data, err := u.datacenterRepo.GetAllUrlCache(pctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
