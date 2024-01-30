package server

import (
	commenthttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentHttpHandler"
	commentusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentUsecase"
)

func (s *server) commentService(commentUsecase *commentusecase.CommentUsecaseService) {
	// commentRepo := commentrepository.NewCommentRepository(s.db)
	// commentUsecase := commentusecase.NewCommentUsecase(commentRepo)
	commentHttpHandler := commenthttphandler.NewCommentHttpHandler(*commentUsecase)

	comments := s.app.Group("/comment_v1")
	comments.POST("/push-comment/:projectId", commentHttpHandler.PushComment, s.middleware.JwtAuthorization)
	comments.PATCH("/update-comment/:projectId/:commentId", commentHttpHandler.UpdateComment, s.middleware.JwtAuthorization)
}
