package projectusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	projectrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"

	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProjectUsecaseService interface {
		CreateNewProjectUsecase(pctx context.Context, req *project.InsertProjectReq) (*project.ProjectRes, error)
		FindOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string) (*project.ProjectResWithUser, error)
		UpdateOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string, req *project.InsertProjectReq) (*project.ProjectResWithUser, error)
		DeleteOneProjectUsecase(pctx context.Context, projectId, userId string) error
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

func (u *projectUsecase) FindOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string) (*project.ProjectResWithUser, error) {

	projectD, err := u.projectRepo.FindOneProject(pctx, utils.ConvertToObjectId(projectId))
	if err != nil {
		return nil, err
	}

	user, err := u.projectRepo.FindOneUserWithId(pctx, grpcCfg.UserUrl, &usersPb.GetUserInfoReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return u.convertPModelToPWithUser(projectD, user)
}

func (u *projectUsecase) UpdateOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string, req *project.InsertProjectReq) (*project.ProjectResWithUser, error) {
	projectD, err := u.projectRepo.UpdateProject(pctx, &project.ProjectModel{
		Id:             utils.ConvertToObjectId(projectId),
		Name:           req.Name,
		LogoUrl:        req.LogoUrl,
		BannerUrl:      req.BannerUrl,
		WebsiteUrl:     req.WebsiteUrl,
		CryptoCategory: req.CryptoCategory,
		Description:    req.Description,
		Reason:         req.Reason,
		Category:       req.Category,
		Contact:        req.Contact,
		UpdatedAt:      utils.LocalTime(),
	}, userId)
	if err != nil {
		return nil, err
	}

	user, err := u.projectRepo.FindOneUserWithId(pctx, grpcCfg.UserUrl, &usersPb.GetUserInfoReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return u.convertPModelToPWithUser(projectD, user)
}

func (u *projectUsecase) DeleteOneProjectUsecase(pctx context.Context, projectId, userId string) error {
	if err := u.projectRepo.DeleteProject(pctx, utils.ConvertToObjectId(projectId), userId); err != nil {
		return err
	}
	return nil
}

func (u *projectUsecase) convertPModelToPWithUser(m *project.ProjectModel, us *usersPb.GetUserInfoRes) (*project.ProjectResWithUser, error) {

	loc, err := utils.LocationTime()
	if err != nil {
		return nil, err
	}
	return &project.ProjectResWithUser{
		Id:             m.Id.Hex(),
		Name:           m.Name,
		LogoUrl:        m.LogoUrl,
		BannerUrl:      m.BannerUrl,
		WebsiteUrl:     m.WebsiteUrl,
		CryptoCategory: m.CryptoCategory,
		Description:    m.Description,
		Reason:         m.Reason,
		Category:       m.Category,
		Contact:        m.Contact,
		FavCount:       m.FavCount,
		CommentCount:   m.CommentCount,
		User: users.UserProfileRes{
			Id:        us.UserId,
			Email:     us.Email,
			Profile:   us.Profile,
			Username:  us.UserName,
			Firstname: us.FirstName,
			Lastname:  us.LastName,
		},
		CreateAt:  m.CreateAt.In(loc),
		UpdatedAt: m.UpdatedAt.In(loc),
	}, nil
}
