package favoritehttphandler

import (
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoritePb"
	favoriteusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteUsecase"
	"context"
)

type (
	favGrpcHandler struct {
		favPb.UnimplementedFavGrpcServiceServer
		favUsecase favoriteusecase.FavoriteUsecaseService
	}
)

func NewfavGrpcHandler(favUsecase favoriteusecase.FavoriteUsecaseService) *favGrpcHandler {
	return &favGrpcHandler{favUsecase: favUsecase}
}

func (g *favGrpcHandler) GetAllFavByUserId(pctx context.Context, req *favPb.GetAllFavProjectReq) (*favPb.GetAllFavRes, error) {
	return g.favUsecase.GetAllProjectByUserId(pctx, req)
}
