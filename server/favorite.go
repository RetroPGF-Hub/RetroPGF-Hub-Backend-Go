package server

import (
	favoritehttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteHttpHandler"
	favoriteusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteUsecase"
)

func (s *server) favoriteService(favoriteUsecase *favoriteusecase.FavoriteUsecaseService) {
	// favoriteRepo := favoriterepository.NewFavoriteRepository(s.db)
	// favoriteUsecase := favoriteusecase.NewFavoriteUsecase(favoriteRepo)
	favoriteHttpHandler := favoritehttphandler.NewFavoriteHttpHandler(*favoriteUsecase)

	favorites := s.app.Group("/fav_v1")
	favorites.POST("/push-pull-fav/:projectId", favoriteHttpHandler.FavPullOrPushHttp, s.middleware.JwtAuthorization)
}
