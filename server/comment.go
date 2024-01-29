package server

import (
	commentPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentPb"

	commenthttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentHttpHandler"
	commentrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentRepository"
	commentusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentUsecase"
	grpcconn "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/grpcConn"
	"log"
)

func (s *server) commentService() {
	commentRepo := commentrepository.NewCommentRepository(s.db)
	commentUsecase := commentusecase.NewCommentUsecase(commentRepo)
	commentHttpHandler := commenthttphandler.NewCommentHttpHandler(commentUsecase)

	commentGrpc := commenthttphandler.NewcommentGrpcHandler(commentUsecase)
	// Grpc client
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.CommentUrl)

		commentPb.RegisterCommentGrpcServiceServer(grpcServer, commentGrpc)

		log.Printf("Comment grpc listening on %s", s.cfg.Grpc.CommentUrl)
		grpcServer.Serve(lis)
	}()

	comments := s.app.Group("/comment_v1")
	comments.POST("/push-comment/:projectId", commentHttpHandler.PushComment, s.middleware.JwtAuthorization)
	comments.PATCH("/update-comment/:projectId/:commentId", commentHttpHandler.UpdateComment, s.middleware.JwtAuthorization)
}
