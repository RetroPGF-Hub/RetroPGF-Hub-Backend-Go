package server

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	favoritehttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteHttpHandler"
)

func (s *server) favoriteService(pActor *modules.ProjectSvcInteractor) {
	// favoriteRepo := favoriterepository.NewFavoriteRepository(s.db)
	// favoriteUsecase := favoriteusecase.NewFavoriteUsecase(favoriteRepo)
	favoriteHttpHandler := favoritehttphandler.NewFavoriteHttpHandler(*pActor)

	// favGrpc := favoritehttphandler.NewfavGrpcHandler(*favoriteUsecase)
	// // Grpc client
	// go func() {
	// 	grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.FavUrl)

	// 	favPb.RegisterFavGrpcServiceServer(grpcServer, favGrpc)

	// 	log.Printf("Favorite grpc listening on %s", s.cfg.Grpc.FavUrl)
	// 	grpcServer.Serve(lis)
	// }()

	favorites := s.app.Group("/fav_v1")
	favorites.POST("/push-pull-fav/:projectId", favoriteHttpHandler.FavPullOrPushHttp, s.middleware.JwtAuthorization)
}
