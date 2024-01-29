package commenthttphandler

import (
	commentPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentPb"
	commentusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentUsecase"
	"context"
	"fmt"
)

type (
	commentGrpcHandler struct {
		commentPb.UnimplementedCommentGrpcServiceServer
		commentUsecase commentusecase.CommentUsecaseService
	}
)

func NewcommentGrpcHandler(commentUsecase commentusecase.CommentUsecaseService) *commentGrpcHandler {
	return &commentGrpcHandler{commentUsecase: commentUsecase}
}

func (g *commentGrpcHandler) CreateCommentProject(ctx context.Context, req *commentPb.CreateCommentProjectReq) (*commentPb.CreateCommentProjectRes, error) {
	if err := g.commentUsecase.InsertEmptyDoc(ctx, req); err != nil {
		return &commentPb.CreateCommentProjectRes{
			Msg: fmt.Sprintf("create comment failed cause of : %s", err.Error()),
		}, err
	}
	return &commentPb.CreateCommentProjectRes{
		Msg: "create comment succes",
	}, nil
}
