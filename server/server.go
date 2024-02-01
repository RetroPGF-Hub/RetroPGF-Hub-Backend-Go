package server

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	commentrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentRepository"
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoritePb"

	commentusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment/commentUsecase"
	favoritehttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteHttpHandler"
	favoriterepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteRepository"
	favoriteusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favoriteUsecase"
	middlewarehttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/middleware/middlewareHttpHandler"
	middlewareusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/middleware/middlewareUsecase"
	projectrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectRepository"
	projectusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project/projectUsecase"
	grpcconn "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/grpcConn"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/jwtauth"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		redis      *redis.Client
		app        *echo.Echo
		db         *mongo.Client
		cfg        *config.Config
		middleware middlewarehttphandler.MiddlewareHttpHandlerService
	}
)

func newMiddleware(cfg *config.Config) middlewarehttphandler.MiddlewareHttpHandlerService {
	usecase := middlewareusecase.NewMiddlewareUsecase()
	return middlewarehttphandler.NewMiddlewareHttpHandler(cfg, usecase)
}

func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {

	log.Printf("Starting service: %s", s.cfg.App.Name)

	<-quit

	log.Printf("Shutting down service: %s", s.cfg.App.Name)

	// depend on which library you use to shutdown the app in this case its fiber
	if err := s.app.Shutdown(pctx); err != nil {
		log.Fatalf("Error: %v", err)
	}

}

func (s *server) httpListening() {
	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
}

func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
	s := &server{
		// redis:      redisactor.RedisConn(&cfg.Redis),
		redis:      nil,
		db:         db,
		cfg:        cfg,
		app:        echo.New(),
		middleware: newMiddleware(cfg),
	}

	jwtauth.SetApiKey(&cfg.Jwt)

	// Request Timeout
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	// CORS
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	// Body Limit
	s.app.Use(middleware.BodyLimit("10M"))
	// app.Settings.MaxRequestBodySize = 10 * 1024 * 1024 // 10 MB

	// Call the server service here
	switch s.cfg.App.Name {
	case "users":
		s.usersService()
	case "datacenter":
		s.datacenterService()
	case "project":

		projectRepo := projectrepository.NewProjectRepository(s.db)
		commentRepo := commentrepository.NewCommentRepository(s.db)
		favoriteRepo := favoriterepository.NewFavoriteRepository(s.db)

		projectActor := modules.NewProjectSvc(projectRepo, commentRepo, favoriteRepo)
		projectUsecase := projectusecase.NewProjectUsecase(*projectActor)

		commentUsecase := commentusecase.NewCommentUsecase(*projectActor)

		favoriteUsecase := favoriteusecase.NewFavoriteUsecase(*projectActor)

		favGrpc := favoritehttphandler.NewfavGrpcHandler(favoriteUsecase)

		// Grpc client
		go func() {
			grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ProjectUrl)

			favPb.RegisterFavGrpcServiceServer(grpcServer, favGrpc)

			log.Printf("Fav grpc listening on %s", s.cfg.Grpc.ProjectUrl)
			grpcServer.Serve(lis)
		}()

		s.projectService(&projectUsecase)

		// comment service
		s.commentService(&commentUsecase)

		// fav service
		s.favoriteService(&favoriteUsecase)
	}

	s.app.Use(middleware.Logger())

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go s.gracefulShutdown(pctx, quit)

	// Listening
	s.httpListening()

}
