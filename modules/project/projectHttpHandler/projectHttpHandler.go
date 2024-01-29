package projecthttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/request"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	ProjectHttpHandlerService interface {
		CreateNewProjectHttp(c echo.Context) error
		FindOneProjectHttp(c echo.Context) error
		UpdateOneProjectHttp(c echo.Context) error
		DeleteOneProjectHttp(c echo.Context) error
	}

	projectHttpHandler struct {
		pActor modules.ProjectSvcInteractor
		cfg    *config.Config
	}
)

func NewProjectHttpHandler(pActor modules.ProjectSvcInteractor, cfg *config.Config) ProjectHttpHandlerService {
	return &projectHttpHandler{
		pActor: pActor,
		cfg:    cfg,
	}
}

func (h *projectHttpHandler) CreateNewProjectHttp(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(project.InsertProjectReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	req.CreatedBy = userId

	res, err := h.pActor.ProjectUsecase.CreateNewProjectUsecase(ctx, req, &h.cfg.Grpc)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":     "ok",
		"project": res,
	})

}

func (h *projectHttpHandler) FindOneProjectHttp(c echo.Context) error {
	ctx := context.Background()

	projectId := c.Param("projectId")
	if len(projectId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "projectId is required")
	}

	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	res, err := h.pActor.ProjectUsecase.FindOneProjectUsecase(ctx, &h.cfg.Grpc, projectId, userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":     "ok",
		"project": res,
	})
}

func (h *projectHttpHandler) DeleteOneProjectHttp(c echo.Context) error {
	ctx := context.Background()

	projectId := c.Param("projectId")

	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	err := h.pActor.ProjectUsecase.DeleteOneProjectUsecase(ctx, projectId, userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg": "delete success",
	})
}

func (h *projectHttpHandler) UpdateOneProjectHttp(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(project.InsertProjectReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	projectId := c.Param("projectId")
	if len(projectId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "projectId is required")
	}

	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	res, err := h.pActor.ProjectUsecase.UpdateOneProjectUsecase(ctx, &h.cfg.Grpc, userId, projectId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":     "ok",
		"project": res,
	})
}
