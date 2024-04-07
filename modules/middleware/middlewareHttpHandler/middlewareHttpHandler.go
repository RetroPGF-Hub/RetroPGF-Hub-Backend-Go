package middlewarehttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	middlewareusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/middleware/middlewareUsecase"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareHttpHandlerService interface {
		JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc
		JwtOptional(next echo.HandlerFunc) echo.HandlerFunc
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

		accessToken, err := c.Cookie("accessToken")
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, "authorization is required")
		}

		if len(accessToken.Value) < 10 {
			return response.ErrResponse(c, http.StatusUnauthorized, "token is required")
		}
		newCtx, err := h.middlewareUsecase.JwtAuthorization(c, h.cfg, accessToken.Value)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}
		return next(newCtx)
	}
}

func (h *middlewareHttpHandler) JwtOptional(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		accessToken := c.Request().Header.Get("Authorization")
		// token exist
		if len(accessToken) > 10 {
			newCtx, err := h.middlewareUsecase.JwtAuthorization(c, h.cfg, accessToken)
			if err != nil {
				return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
			}
			return next(newCtx)

			// token doesn't exist
		} else {
			// fmt.Println("token doesn't exist")
			subAccessToken, err := c.Cookie("accessToken")
			if errors.Is(err, http.ErrNoCookie) {
				c.Set("user_id", "")
				c.Set("email", "")
				c.Set("source", "")
				return next(c)
			}

			if len(subAccessToken.Value) > 10 {
				// fmt.Println("token exist")
				newCtx, err := h.middlewareUsecase.JwtAuthorization(c, h.cfg, subAccessToken.Value)
				if err != nil {
					return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
				}
				return next(newCtx)
			} else {
				c.Set("user_id", "")
				c.Set("email", "")
				c.Set("source", "")
				return next(c)
			}

		}

	}
}
