package commentusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CommentUsecaseService interface {
		PushCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId string) (*comment.CommentARes, error)
		UpdateCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId, commentId string) (*comment.CommentProjectRes, error)
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

func (u *commentUsecase) PushCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId string) (*comment.CommentARes, error) {

	projectIdPri := utils.ConvertToObjectId(projectId)
	countProject, err := u.pActor.CommentRepo.CountComment(pctx, projectIdPri)
	if err != nil {
		return nil, err
	}

	if countProject == 0 {
		return nil, errors.New("project you looking for is not exist")
	}

	commentId := primitive.NewObjectID()
	commentA := new(comment.CommentA)
	commentA.CommentId = commentId
	commentA.Title = req.Title
	commentA.Content = req.Content
	commentA.CreatedBy = req.CreatedBy
	commentA.CreateAt = utils.LocalTime()
	commentA.UpdatedAt = utils.LocalTime()

	if err := u.pActor.CommentRepo.PushComment(pctx, projectIdPri, commentA); err != nil {
		return nil, err
	}

	if err := u.pActor.ProjectRepo.UpdateProjectCommentCount(pctx, projectIdPri, 1); err != nil {
		return nil, err
	}

	return &comment.CommentARes{
		CommentId: commentA.CommentId.Hex(),
		Title:     commentA.Title,
		Content:   commentA.Content,
		CreatedBy: commentA.CreatedBy,
		CreateAt:  commentA.CreateAt,
		UpdatedAt: commentA.UpdatedAt,
	}, nil

}

func (u *commentUsecase) UpdateCommentUsecase(pctx context.Context, req *comment.PushCommentReq, projectId, commentId string) (*comment.CommentProjectRes, error) {

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

	var CommentProjectRes []comment.CommentARes

	for _, v := range comments.Comments {
		CommentProjectRes = append(CommentProjectRes, comment.CommentARes{
			CommentId: v.CommentId.Hex(),
			Content:   v.Content,
			Title:     v.Title,
			CreatedBy: v.CreatedBy,
			CreateAt:  v.CreateAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	fmt.Println("data ->", CommentProjectRes)
	return &comment.CommentProjectRes{
		ProjectId: comments.ProjectId.Hex(),
		Comments:  CommentProjectRes,
		CreateAt:  comments.CreateAt,
		UpdatedAt: comments.UpdatedAt,
	}, nil

}
