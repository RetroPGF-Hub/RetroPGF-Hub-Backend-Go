package usersusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoritePb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
	usersrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/jwtauth"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type (
	UsersUsecaseService interface {
		RegisterUserUsecase(cfg *config.Config, pctx context.Context, req *users.RegisterUserReq) (string, *users.UserProfileRes, error)
		LoginUsecase(cfg *config.Config, pctx context.Context, email, password string) (string, *users.UserProfileRes, error)
		FindUserByIdUsecase(pctx context.Context, req *usersPb.GetUserInfoReq) (*usersPb.GetUserInfoRes, error)
		GetUserFavs(pctx context.Context, cfg *config.Grpc, userId string) (*favPb.GetAllFavRes, error)
	}

	usersUsecase struct {
		usersRepo usersrepository.UsersRepositoryService
	}
)

func NewUsersUsecase(usersRepo usersrepository.UsersRepositoryService) UsersUsecaseService {
	return &usersUsecase{
		usersRepo: usersRepo,
	}
}

func (u *usersUsecase) RegisterUserUsecase(cfg *config.Config, pctx context.Context, req *users.RegisterUserReq) (string, *users.UserProfileRes, error) {
	exist, err := u.usersRepo.IsUniqueUser(pctx, req.Email)
	if err != nil {
		return "", nil, err
	}

	if exist {
		return "", nil, errors.New("email is already exist, try to use different email")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, errors.New("error: failed to hash password")
	}

	userId, err := u.usersRepo.InsertOneUser(pctx, &users.UserDb{
		Id:        primitive.NewObjectID(),
		Email:     req.Email,
		Password:  string(hashPassword),
		Source:    req.Source,
		Profile:   req.Profile,
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		CreateAt:  utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
	})

	if err != nil {
		return "", nil, err
	}

	user, err := u.usersRepo.FindOneUserWithIdWithPassword(pctx, userId)
	if err != nil {
		return "", nil, err
	}

	accessToken, err := jwtauth.NewAccessToken(
		&cfg.Jwt,
		&jwtauth.Claims{
			UserId:   user.Id.Hex(),
			Email:    user.Email,
			Username: user.Username,
			Source:   user.Source,
			Profile:  user.Profile,
		}, cfg.Jwt.AccessDuration, "accessToken").SignToken()
	if err != nil {
		return "", nil, err
	}

	return accessToken,
		&users.UserProfileRes{
			Email:     user.Email,
			Profile:   user.Profile,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Id:        user.Id.Hex(),
		},
		nil

}

func (u *usersUsecase) LoginUsecase(cfg *config.Config, pctx context.Context, email, password string) (string, *users.UserProfileRes, error) {
	result, err := u.usersRepo.FindOneUserWithEmail(pctx, email)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		return "", nil, errors.New("error: email or password invalid")
	}

	accessToken, err := jwtauth.NewAccessToken(
		&cfg.Jwt,
		&jwtauth.Claims{
			UserId:   result.Id.Hex(),
			Email:    result.Email,
			Username: result.Username,
			Source:   result.Source,
			Profile:  result.Profile,
		}, cfg.Jwt.AccessDuration, "accessToken").SignToken()
	if err != nil {
		return "", nil, err
	}

	return accessToken,
		&users.UserProfileRes{
			Email:     result.Email,
			Profile:   result.Profile,
			Username:  result.Username,
			Firstname: result.Firstname,
			Lastname:  result.Lastname,
			Id:        result.Id.Hex(),
		},
		nil
}

func (u *usersUsecase) GetUserFavs(pctx context.Context, cfg *config.Grpc, userId string) (*favPb.GetAllFavRes, error) {
	projects, err := u.usersRepo.GetFavProjectByUserId(pctx, cfg.FavUrl, userId)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
