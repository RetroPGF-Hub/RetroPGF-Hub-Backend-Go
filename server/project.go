package server

import (
	projecthttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectHttpHandler"
	projectrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectRepository"
	projectusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectUsecase"
)

func (s *server) projectService() {
	projectRepo := projectrepository.NewProjectRepository(s.db)
	projectUsecase := projectusecase.NewProjectUsecase(projectRepo)
	projectHttpHandler := projecthttphandler.NewProjectHttpHandler(projectUsecase, s.cfg)

	projects := s.app.Group("/project_v1")
	projects.POST("/create", projectHttpHandler.CreateNewProjectHttp, s.middleware.JwtAuthorization)
	projects.GET("/project/:projectId", projectHttpHandler.FindOneProjectHttp, s.middleware.JwtAuthorization)
	projects.PATCH("/project/:projectId", projectHttpHandler.UpdateOneProjectHttp, s.middleware.JwtAuthorization)
	projects.DELETE("/project/:projectId", projectHttpHandler.DeleteOneProjectHttp, s.middleware.JwtAuthorization)
}
