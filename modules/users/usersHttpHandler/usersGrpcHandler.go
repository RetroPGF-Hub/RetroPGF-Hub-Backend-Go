package usershttphandler

import (
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
	usersusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersUsecase"
	"context"
)

type (
	usersGrpcHandler struct {
		usersPb.UnimplementedUsersGrpcServiceServer
		usersUsecase usersusecase.UsersUsecaseService
	}
)

func NewusersGrpcHandler(usersUsecase usersusecase.UsersUsecaseService) *usersGrpcHandler {
	return &usersGrpcHandler{usersUsecase: usersUsecase}
}

func (g *usersGrpcHandler) GetUserInfoById(ctx context.Context, req *usersPb.GetUserInfoReq) (*usersPb.GetUserInfoRes, error) {
	return g.usersUsecase.FindUserByIdUsecase(ctx, req)
}
