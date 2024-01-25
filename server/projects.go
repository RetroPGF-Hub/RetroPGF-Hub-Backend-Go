package server

func (s *server) projectsService() {
	projectsRepo := projectsrepository.NewProjectsRepository(s.db)
	projectsUsecase := projectsusecase.NewProjectssUsecase(projectsRepo)
	projectsHttpHandler := projectshandler.NewProjectssHttpHandler(s.cfg, projectsUsecase)
	
	_ = projectsHttpHandler
}

