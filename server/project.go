package server

import (
	projecthttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectHttpHandler"
	projectusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectUsecase"
)

func (s *server) projectService(projectUsecase *projectusecase.ProjectUsecaseService) {
	// projectRepo := projectrepository.NewProjectRepository(s.db)
	// projectUsecase := projectusecase.NewProjectUsecase(projectRepo)
	projectHttpHandler := projecthttphandler.NewProjectHttpHandler(*projectUsecase, s.cfg)

	projects := s.app.Group("/project_v1")
	projects.POST("/create", projectHttpHandler.CreateNewProjectHttp, s.middleware.JwtAuthorization)
	projects.GET("/projects", projectHttpHandler.FindAllProeject, s.middleware.JwtOptional)
	projects.GET("/project/:projectId", projectHttpHandler.FindOneProjectHttp, s.middleware.JwtOptional)
	projects.PATCH("/project/:projectId", projectHttpHandler.UpdateOneProjectHttp, s.middleware.JwtAuthorization)
	projects.DELETE("/project/:projectId", projectHttpHandler.DeleteOneProjectHttp, s.middleware.JwtAuthorization)

}
