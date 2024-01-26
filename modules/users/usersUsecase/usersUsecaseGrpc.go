package usersusecase

import (
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
)

func (u *usersUsecase) FindUserByIdUsecase(pctx context.Context, req *usersPb.GetUserInfoReq) (*usersPb.GetUserInfoRes, error) {
	user, err := u.usersRepo.FindOneUserWithId(pctx, utils.ConvertToObjectId(req.UserId))
	if err != nil {
		return nil, err
	}

	return &usersPb.GetUserInfoRes{
		UserId:    user.Id.Hex(),
		Email:     user.Email,
		Source:    user.Source,
		Profile:   user.Profile,
		UserName:  user.Username,
		FirstName: user.Firstname,
		LastName:  user.Lastname,
	}, nil

}
