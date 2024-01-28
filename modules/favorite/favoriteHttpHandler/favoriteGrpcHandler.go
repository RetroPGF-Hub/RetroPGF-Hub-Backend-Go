package favoritehttphandler

import (
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favPb"
	favoriteusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteUsecase"
	"context"
	"fmt"
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

func (g *favGrpcHandler) CreateFavProject(ctx context.Context, req *favPb.CreateFavProjectReq) (*favPb.CreateFavProjectRes, error) {
	if err := g.favUsecase.InsertEmptyDoc(ctx, req); err != nil {
		return &favPb.CreateFavProjectRes{
			Msg: fmt.Sprintf("create fav failed cause of : %s", err.Error()),
		}, err
	}
	return &favPb.CreateFavProjectRes{
		Msg: "create fav succes",
	}, nil
}
