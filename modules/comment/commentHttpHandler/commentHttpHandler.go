package commenthttphandler

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/request"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/response"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	CommentHttpHandlerService interface {
		PushComment(c echo.Context) error
		UpdateComment(c echo.Context) error
	}

	commentHttpHandler struct {
		pActor modules.ProjectSvcInteractor
	}
)

func NewCommentHttpHandler(pActor modules.ProjectSvcInteractor) CommentHttpHandlerService {
	return &commentHttpHandler{
		pActor: pActor,
	}
}

func (h *commentHttpHandler) PushComment(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(comment.PushCommentReq)

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

	req.CreatedBy = userId

	err := h.pActor.CommentUsecase.PushCommentUsecase(ctx, req, projectId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg": "create comment success",
	})

}

func (h *commentHttpHandler) UpdateComment(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(comment.PushCommentReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	projectId := c.Param("projectId")
	if len(projectId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "projectId is required")
	}

	commentId := c.Param("commentId")
	if len(commentId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "commentId is required")
	}

	userId := c.Get("user_id").(string)
	if len(userId) < 5 {
		return response.ErrResponse(c, http.StatusBadRequest, "unauthorized user")
	}

	req.CreatedBy = userId

	updatedComment, err := h.pActor.CommentUsecase.UpdateCommentUsecase(ctx, req, projectId, commentId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"msg":     "update comment success",
		"comment": updatedComment,
	})

}
