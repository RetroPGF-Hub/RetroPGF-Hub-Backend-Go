package projectusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	projectrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProjectUsecaseService interface {
		CreateNewProjectUsecase(pctx context.Context, req *project.InsertProjectReq) (*project.ProjectRes, error)
	}

	projectUsecase struct {
		projectRepo projectrepository.ProjectRepositoryService
	}
)

func NewProjectUsecase(projectRepo projectrepository.ProjectRepositoryService) ProjectUsecaseService {
	return &projectUsecase{
		projectRepo: projectRepo,
	}
}

func (u *projectUsecase) CreateNewProjectUsecase(pctx context.Context, req *project.InsertProjectReq) (*project.ProjectRes, error) {
	projectId, err := u.projectRepo.InsertOneProject(pctx, &project.ProjectModel{
		Id:             primitive.NewObjectID(),
		Name:           req.Name,
		LogoUrl:        req.LogoUrl,
		FavCount:       0,
		CommentCount:   0,
		BannerUrl:      req.BannerUrl,
		WebsiteUrl:     req.WebsiteUrl,
		CryptoCategory: req.CryptoCategory,
		Description:    req.Description,
		Reason:         req.Reason,
		Category:       req.Category,
		Contact:        req.Contact,
		CreatedBy:      req.CreatedBy,
		CreateAt:       utils.LocalTime(),
		UpdatedAt:      utils.LocalTime(),
	})
	if err != nil {
		return nil, err
	}

	projectData, err := u.projectRepo.FindOneProject(pctx, projectId)
	if err != nil {
		return nil, err
	}

	return &project.ProjectRes{
		Id:             projectData.Id.Hex(),
		Name:           projectData.Name,
		LogoUrl:        projectData.LogoUrl,
		BannerUrl:      projectData.BannerUrl,
		WebsiteUrl:     projectData.WebsiteUrl,
		CryptoCategory: projectData.CryptoCategory,
		Description:    projectData.Description,
		Reason:         projectData.Reason,
		Category:       projectData.Category,
		FavCount:       projectData.FavCount,
		CommentCount:   projectData.CommentCount,
		Contact:        projectData.Contact,
		CreatedBy:      projectData.CreatedBy,
		CreateAt:       projectData.CreateAt,
		UpdatedAt:      projectData.UpdatedAt,
	}, nil
}

func (u *projectUsecase) FindOneProjectUsecase(pctx context.Context, projectId string) (*project.ProjectRes, error) {

	projectData, err := u.projectRepo.FindOneProject(pctx, utils.ConvertToObjectId(projectId))
	if err != nil {
		return nil, err
	}
	_ = projectData

	return &project.ProjectRes{}, nil
}
