package usersusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	usersrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UsersUsecaseService interface {
		InsertOneUser(pctx context.Context, req *users.RegisterUserReq) (*users.UserProfileRes, error)
	}

	usersUsecase struct {
		usersRepo usersrepository.UsersRepositoryService
	}
)

func NewUsersRepo(usersRepo usersrepository.UsersRepositoryService) UsersUsecaseService {
	return &usersUsecase{
		usersRepo: usersRepo,
	}
}

func (u *usersUsecase) InsertOneUser(pctx context.Context, req *users.RegisterUserReq) (*users.UserProfileRes, error) {
	exist, err := u.usersRepo.IsUniqueUser(pctx, req.Email)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errors.New("email is already exist, try to use different email")
	}

	userId, err := u.usersRepo.InsertOneUser(pctx, &users.UserDb{
		Id:        primitive.NewObjectID(),
		Email:     req.Email,
		Profile:   req.Profile,
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		CreateAt:  utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
	})

	if err != nil {
		return nil, err
	}

	user, err := u.usersRepo.FindOneUserWithId(pctx, userId)
	if err != nil {
		return nil, err
	}

	return &users.UserProfileRes{
		Email:     user.Email,
		Profile:   user.Profile,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Id:        user.Id.Hex(),
	}, nil

}
