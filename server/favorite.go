package server

import (
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favPb"
	favoritehttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteHttpHandler"
	favoriterepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteRepository"
	favoriteusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteUsecase"
	grpcconn "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/grpcConn"
	"log"
)

func (s *server) favoriteService() {
	favoriteRepo := favoriterepository.NewFavoriteRepository(s.db)
	favoriteUsecase := favoriteusecase.NewFavoriteUsecase(favoriteRepo)
	favoriteHttpHandler := favoritehttphandler.NewFavoriteHttpHandler(favoriteUsecase)

	favGrpc := favoritehttphandler.NewfavGrpcHandler(favoriteUsecase)
	// Grpc client
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.FavUrl)

		favPb.RegisterFavGrpcServiceServer(grpcServer, favGrpc)

		log.Printf("Favorite grpc listening on %s", s.cfg.Grpc.FavUrl)
		grpcServer.Serve(lis)
	}()

	favorites := s.app.Group("/fav_v1")
	favorites.POST("/push-pull-fav/:projectId", favoriteHttpHandler.FavPullOrPushHttp, s.middleware.JwtAuthorization)
}
