package server

import (
	projecthttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectHttpHandler"
	projectrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectRepository"
	projectusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectUsecase"
)

func (s *server) projectService() {
	projectRepo := projectrepository.NewProjectRepository(s.db)
	projectUsecase := projectusecase.NewProjectRepo(projectRepo)
	projectHttpHandler := projecthttphandler.NewProjectHttpHandler(projectUsecase)

	_ = projectHttpHandler
}
