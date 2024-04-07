package middlewareusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/jwtauth"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error)
	}

	middlewareUsecase struct {
	}
)

func NewMiddlewareUsecase() MiddlewareUsecaseService {
	return &middlewareUsecase{}

}

func (u *middlewareUsecase) JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error) {

	cliams, err := jwtauth.ParseToken(accessToken, &cfg.Jwt)
	if err != nil {
		return nil, err
	}

	// fmt.Println("user id set ->", cliams.UserId)
	c.Set("user_id", cliams.UserId)
	c.Set("email", cliams.Email)
	c.Set("source", cliams.Source)

	return c, nil
}
