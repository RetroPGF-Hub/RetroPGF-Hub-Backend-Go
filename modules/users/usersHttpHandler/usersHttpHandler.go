package usershttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	usersusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersUsecase"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/request"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	UsersHttpHandlerService interface {
		RegisterUser(c echo.Context) error
		LoginUser(c echo.Context) error
		LogOutUser(c echo.Context) error
	}

	usersHttpHandler struct {
		cfg          *config.Config
		usersUsecase usersusecase.UsersUsecaseService
	}
)

func NewUsersHttpHandler(cfg *config.Config, usersUsecase usersusecase.UsersUsecaseService) UsersHttpHandlerService {
	return &usersHttpHandler{
		usersUsecase: usersUsecase,
		cfg:          cfg,
	}
}

func (h *usersHttpHandler) RegisterUser(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(users.RegisterUserReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	token, res, err := h.usersUsecase.RegisterUserUsecase(h.cfg, ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	oneWeek := time.Now().Add(7 * 24 * time.Hour)

	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Expires:  oneWeek,
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessChecker",
		Value:    res.Id,
		HttpOnly: false,
		Secure:   false,
		Expires:  oneWeek,
	})

	return response.SuccessResponse(c, http.StatusCreated, map[string]any{
		"msg":  "ok",
		"user": res,
	})

}

func (h *usersHttpHandler) LoginUser(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(users.LoginReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	token, res, err := h.usersUsecase.LoginUsecase(h.cfg, ctx, req.Email, req.Password)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	oneWeek := time.Now().Add(7 * 24 * time.Hour)

	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Expires:  oneWeek,
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessChecker",
		Value:    res.Id,
		HttpOnly: false,
		Secure:   false,
		Expires:  oneWeek,
	})

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":  "ok",
		"user": res,
	})
}

func (h *usersHttpHandler) LogOutUser(c echo.Context) error {

	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessChecker",
		Value:    "",
		HttpOnly: false,
		Secure:   false,
		Expires:  time.Unix(0, 0),
	})

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg": "logout success",
	})

}
