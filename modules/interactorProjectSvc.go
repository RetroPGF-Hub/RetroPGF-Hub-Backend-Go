package modules

import (
	commentusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentUsecase"
	favoriteusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteUsecase"
	projectusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectUsecase"
)

type ProjectSvcInteractor struct {
	ProjectUsecase  projectusecase.ProjectUsecaseService
	CommentUsecase  commentusecase.CommentUsecaseService
	FavoriteUsecase favoriteusecase.FavoriteUsecaseService
}

func NewProjectSvc(projectUsecase projectusecase.ProjectUsecaseService, commentUsecase commentusecase.CommentUsecaseService, favoriteUsecase favoriteusecase.FavoriteUsecaseService) *ProjectSvcInteractor {
	return &ProjectSvcInteractor{
		ProjectUsecase:  projectUsecase,
		CommentUsecase:  commentUsecase,
		FavoriteUsecase: favoriteUsecase,
	}
}
