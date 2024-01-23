package projecthttphandler

import projectusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectUsecase"

type (
	ProjectHttpHandlerService interface {
	}

	projectHttpHandler struct {
		projectUsecase projectusecase.ProjectUsecaseService
	}
)

func NewProjectHttpHandler(projectUsecase projectusecase.ProjectUsecaseService) ProjectHttpHandlerService {
	return &projectHttpHandler{
		projectUsecase: projectUsecase,
	}
}
