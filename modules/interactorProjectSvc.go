package modules

import (
	commentrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentRepository"
	favoriterepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteRepository"
	projectrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectRepository"
)

type ProjectSvcInteractor struct {
	ProjectRepo  projectrepository.ProjectRepositoryService
	CommentRepo  commentrepository.CommentRepositoryService
	FavoriteRepo favoriterepository.FavoriteRepositoryService
}

func NewProjectSvc(projectRepo projectrepository.ProjectRepositoryService, commentRepo commentrepository.CommentRepositoryService, favoriteRepo favoriterepository.FavoriteRepositoryService) *ProjectSvcInteractor {
	return &ProjectSvcInteractor{
		ProjectRepo:  projectRepo,
		CommentRepo:  commentRepo,
		FavoriteRepo: favoriteRepo,
	}
}
