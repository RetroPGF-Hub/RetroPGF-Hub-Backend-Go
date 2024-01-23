package usershttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	usersusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersUsecase"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/request"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	UsersHttpHandlerService interface {
	}

	usersHttpHandler struct {
		usersUsecase usersusecase.UsersUsecaseService
	}
)

func NewUsersHttpHandler(usersUsecase usersusecase.UsersUsecaseService) UsersHttpHandlerService {
	return &usersHttpHandler{
		usersUsecase: usersUsecase,
	}
}

func (h *usersHttpHandler) RegisterUser(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(users.RegisterUserReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.usersUsecase.InsertOneUser(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, map[string]any{
		"msg":  "ok",
		"user": res,
	})

}
