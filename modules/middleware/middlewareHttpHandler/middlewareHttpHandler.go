package middlewarehttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	middlewareusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/middleware/middlewareUsecase"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareHttpHandlerService interface {
		JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc
	}

	middlewareHttpHandler struct {
		cfg               *config.Config
		middlewareUsecase middlewareusecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHttpHandler(cfg *config.Config, middlewareUsecase middlewareusecase.MiddlewareUsecaseService) MiddlewareHttpHandlerService {
	return &middlewareHttpHandler{
		middlewareUsecase: middlewareUsecase,
		cfg:               cfg,
	}
}

func (h *middlewareHttpHandler) JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		accessToken := c.Request().Header.Get("accessToken")
		if len(accessToken) < 10 {
			return response.ErrResponse(c, http.StatusUnauthorized, "token is required")
		}
		newCtx, err := h.middlewareUsecase.JwtAuthorization(c, h.cfg, accessToken)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}
		return next(newCtx)
	}
}
