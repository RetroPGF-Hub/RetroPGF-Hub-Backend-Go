package server

import (
	usershttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersHttpHandler"
	usersrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersRepository"
	usersusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersUsecase"
)

func (s *server) usersService() {
	usersRepo := usersrepository.NewUsersRepository(s.db)
	usersUsecase := usersusecase.NewUsersUsecase(usersRepo)
	usersHttpHandler := usershttphandler.NewUsersHttpHandler(s.cfg, usersUsecase)

	users := s.app.Group("/users_v1")
	users.POST("/register", usersHttpHandler.RegisterUser)
	users.POST("/login", usersHttpHandler.LoginUser)
	users.GET("/logout", usersHttpHandler.LogOutUser)
}
