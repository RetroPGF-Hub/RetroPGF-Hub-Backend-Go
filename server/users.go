package server

import (
	usershttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersHttpHandler"
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
	usersrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersRepository"
	usersusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersUsecase"
	grpcconn "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/grpcConn"
	"log"
)

func (s *server) usersService() {
	usersRepo := usersrepository.NewUsersRepository(s.db)
	usersUsecase := usersusecase.NewUsersUsecase(usersRepo)
	usersHttpHandler := usershttphandler.NewUsersHttpHandler(s.cfg, usersUsecase)

	userGrpc := usershttphandler.NewusersGrpcHandler(usersUsecase)
	// Grpc client
	go func() {
		log.Println(s.cfg.Grpc.UserUrl)
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.UserUrl)

		usersPb.RegisterUsersGrpcServiceServer(grpcServer, userGrpc)

		log.Printf("Auth grpc listening on %s", s.cfg.Grpc.UserUrl)
		grpcServer.Serve(lis)
	}()

	users := s.app.Group("/users_v1")
	users.POST("/register", usersHttpHandler.RegisterUser)
	users.POST("/login", usersHttpHandler.LoginUser)
	users.GET("/logout", usersHttpHandler.LogOutUser)
}
