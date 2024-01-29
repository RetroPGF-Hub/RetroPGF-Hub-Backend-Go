package server

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	commenthttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentHttpHandler"
)

func (s *server) commentService(pActor *modules.ProjectSvcInteractor) {
	// commentRepo := commentrepository.NewCommentRepository(s.db)
	// commentUsecase := commentusecase.NewCommentUsecase(commentRepo)
	commentHttpHandler := commenthttphandler.NewCommentHttpHandler(*pActor)

	// commentGrpc := commenthttphandler.NewcommentGrpcHandler(*commentUsecase)
	// // Grpc client
	// go func() {
	// 	grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.CommentUrl)

	// 	commentPb.RegisterCommentGrpcServiceServer(grpcServer, commentGrpc)

	// 	log.Printf("Comment grpc listening on %s", s.cfg.Grpc.CommentUrl)
	// 	grpcServer.Serve(lis)
	// }()

	comments := s.app.Group("/comment_v1")
	comments.POST("/push-comment/:projectId", commentHttpHandler.PushComment, s.middleware.JwtAuthorization)
	comments.PATCH("/update-comment/:projectId/:commentId", commentHttpHandler.UpdateComment, s.middleware.JwtAuthorization)
}
