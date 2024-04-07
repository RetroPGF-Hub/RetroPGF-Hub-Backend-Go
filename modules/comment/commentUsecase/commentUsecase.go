package commentusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"

	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CommentUsecaseService interface {
		PushCommentUsecase(grpcCfg *config.Grpc, pctx context.Context, req *comment.PushCommentReq, projectId string) (*comment.CommentAResWithUser, error)
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

func (u *commentUsecase) PushCommentUsecase(grpcCfg *config.Grpc, pctx context.Context, req *comment.PushCommentReq, projectId string) (*comment.CommentAResWithUser, error) {

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
	commentA.Content = req.Content
	commentA.CreatedBy = utils.ConvertToObjectId(req.CreatedBy)
	commentA.CreateAt = utils.LocalTime()
	commentA.UpdatedAt = utils.LocalTime()

	if err := u.pActor.CommentRepo.PushComment(pctx, projectIdPri, commentA); err != nil {
		return nil, err
	}

	if err := u.pActor.ProjectRepo.UpdateProjectCommentCount(pctx, projectIdPri, 1); err != nil {
		return nil, err
	}

	// u.pActor.ProjectRepo.FindOneUserWithId(pctx, )
	user, err := u.pActor.ProjectRepo.FindOneUserWithId(pctx, grpcCfg.UserUrl, &usersPb.GetUserInfoReq{
		UserId: req.CreatedBy,
	})
	if err != nil {
		return nil, err
	}

	return &comment.CommentAResWithUser{
		CommentId: commentA.CommentId.Hex(),
		Content:   commentA.Content,
		CreatedBy: users.SecureUserProfile{
			Profile:  user.Profile,
			Username: user.UserName,
		},
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
		Content:   req.Content,
		CreatedBy: utils.ConvertToObjectId(req.CreatedBy),
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
			CreatedBy: v.CreatedBy.Hex(),
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
