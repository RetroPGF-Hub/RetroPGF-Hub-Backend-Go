package projecthttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	projectusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectUsecase"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/request"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	ProjectHttpHandlerService interface {
		CreateNewProjectHttp(c echo.Context) error
	}

	projectHttpHandler struct {
		projectUsecase projectusecase.ProjectUsecaseService
	}
)

func NewProjectHttpHandler(projectUsecase projectusecase.ProjectUsecaseService) ProjectHttpHandlerService {
	return &projectHttpHandler{
		projectUsecase: projectUsecase,
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
	req.CreatedBy = userId
	// email := c.Get("email")
	// source := c.Get("source")
	// fmt.Println(userId, email, source)

	res, err := h.projectUsecase.CreateNewProjectUsecase(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":     "ok",
		"project": res,
	})

}
