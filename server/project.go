package server

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	projecthttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectHttpHandler"
)

func (s *server) projectService(pActor *modules.ProjectSvcInteractor) {
	// projectRepo := projectrepository.NewProjectRepository(s.db)
	// projectUsecase := projectusecase.NewProjectUsecase(projectRepo)
	projectHttpHandler := projecthttphandler.NewProjectHttpHandler(*pActor, s.cfg)

	projects := s.app.Group("/project_v1")
	projects.POST("/create", projectHttpHandler.CreateNewProjectHttp, s.middleware.JwtAuthorization)
	projects.GET("/project/:projectId", projectHttpHandler.FindOneProjectHttp, s.middleware.JwtAuthorization)
	projects.PATCH("/project/:projectId", projectHttpHandler.UpdateOneProjectHttp, s.middleware.JwtAuthorization)
	projects.DELETE("/project/:projectId", projectHttpHandler.DeleteOneProjectHttp, s.middleware.JwtAuthorization)
}
