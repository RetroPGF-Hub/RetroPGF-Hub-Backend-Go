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
		RegisterOrLogin(c echo.Context) error
		LogOutUser(c echo.Context) error
		GetUserFav(c echo.Context) error
		GetCurrentUser(c echo.Context) error
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
		Expires:  oneWeek,
		Path:     "/",
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessChecker",
		Value:    res.Id,
		Expires:  oneWeek,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
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
		Expires:  oneWeek,
		Path:     "/",
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessChecker",
		Value:    res.Id,
		Expires:  oneWeek,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":  "ok",
		"user": res,
	})
}

func (h *usersHttpHandler) RegisterOrLogin(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(users.RegisterUserReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	token, res, err := h.usersUsecase.RegisterOrLoginThridParty(h.cfg, ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	oneWeek := time.Now().Add(7 * 24 * time.Hour)
	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    token,
		HttpOnly: true,
		Expires:  oneWeek,
		Path:     "/",
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessChecker",
		Value:    res.Id,
		Expires:  oneWeek,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
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
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessChecker",
		Value:    "",
		HttpOnly: false,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg": "logout success",
	})

}

func (h *usersHttpHandler) GetUserFav(c echo.Context) error {

	ctx := context.Background()

	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	projects, err := h.usersUsecase.GetUserFavs(ctx, &h.cfg.Grpc, userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":  "ok",
		"favs": projects,
	})

}

func (h *usersHttpHandler) GetCurrentUser(c echo.Context) error {
	ctx := context.Background()
	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	user, err := h.usersUsecase.GetCurrentUserUsecase(ctx, userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":  "ok",
		"user": user,
	})
}
