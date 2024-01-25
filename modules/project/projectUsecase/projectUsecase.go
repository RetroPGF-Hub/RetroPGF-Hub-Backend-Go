package projectusecase

import projectrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectRepository"

type (
	ProjectUsecaseService interface {
	}

	projectUsecase struct {
		projectRepo projectrepository.ProjectRepositoryService
	}
)

func NewProjectRepo(projectRepo projectrepository.ProjectRepositoryService) ProjectUsecaseService {
	return &projectUsecase{
		projectRepo: projectRepo,
	}
}

// func (u *projectUsecase) CreateNewProject(pctx context.Context)
