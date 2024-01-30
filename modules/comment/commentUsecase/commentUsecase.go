package commentusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CommentUsecaseService interface {
		PushCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId string) error
		UpdateCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId, commentId string) (*comment.CommentRes, error)
	}

	commentUsecase struct {
		pActor modules.ProjectSvcInteractor
	}
)

func NewCommentUsecase(pActor modules.ProjectSvcInteractor) CommentUsecaseService {
	return &commentUsecase{
		pActor: pActor,
	}
}

func (u *commentUsecase) PushCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId string) error {

	projectIdPri := utils.ConvertToObjectId(projectId)
	countProject, err := u.pActor.CommentRepo.CountComment(pctx, projectIdPri)
	if err != nil {
		return err
	}

	if countProject == 0 {
		return errors.New("project you looking for is not exist")
	}

	if err := u.pActor.CommentRepo.PushComment(pctx, projectIdPri, &comment.CommentA{
		CommentId: primitive.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		CreatedBy: req.CreatedBy,
		CreateAt:  utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
	}); err != nil {
		return err
	}

	if err := u.pActor.ProjectRepo.UpdateCommentCount(pctx, projectIdPri, 1); err != nil {

		return err
	}

	return nil

}

func (u *commentUsecase) UpdateCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId, commentId string) (*comment.CommentRes, error) {

	projectIdPri := utils.ConvertToObjectId(projectId)
	countProject, err := u.pActor.CommentRepo.CountComment(pctx, projectIdPri)
	if err != nil {
		return nil, err
	}

	if countProject == 0 {
		return nil, errors.New("project you looking for is not exist")
	}

	comments, err := u.pActor.CommentRepo.UpdateComment(pctx, projectIdPri, &comment.CommentA{
		CommentId: utils.ConvertToObjectId(commentId),
		Title:     req.Title,
		Content:   req.Content,
		CreatedBy: req.CreatedBy,
		CreateAt:  utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
	})
	if err != nil {
		return nil, err
	}

	var commentRes []comment.CommentARes

	for _, v := range comments.Comments {
		commentRes = append(commentRes, comment.CommentARes{
			CommentId: v.CommentId.Hex(),
			Content:   v.Content,
			Title:     v.Title,
			CreatedBy: v.CreatedBy,
			CreateAt:  v.CreateAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return &comment.CommentRes{
		ProjectId: comments.ProjectId.Hex(),
		Comments:  commentRes,
		CreateAt:  comments.CreateAt,
		UpdatedAt: comments.UpdatedAt,
	}, nil

}
