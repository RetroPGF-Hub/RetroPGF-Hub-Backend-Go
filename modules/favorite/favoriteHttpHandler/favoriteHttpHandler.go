package favoritehttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	FavoriteHttpHandlerService interface {
		FavPullOrPushHttp(c echo.Context) error
	}

	favoriteHttpHandler struct {
		pActor modules.ProjectSvcInteractor
	}
)

func NewFavoriteHttpHandler(pActor modules.ProjectSvcInteractor) FavoriteHttpHandlerService {
	return &favoriteHttpHandler{
		pActor: pActor,
	}
}

func (h *favoriteHttpHandler) FavPullOrPushHttp(c echo.Context) error {
	ctx := context.Background()
	projectId := c.Param("projectId")
	if len(projectId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "projectId is required")
	}
	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	opera, err := h.pActor.FavoriteUsecase.FavPullOrPushUsecase(ctx, projectId, userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]string{
		"msg":   "success",
		"opera": opera,
	})
}
