package server

func (s *server) usersService() {
	usersRepo := usersrepository.NewUsersRepository(s.db)
	usersUsecase := usersusecase.NewUserssUsecase(usersRepo)
	usersHttpHandler := usershandler.NewUserssHttpHandler(s.cfg, usersUsecase)
	
	_ = usersHttpHandler
}

