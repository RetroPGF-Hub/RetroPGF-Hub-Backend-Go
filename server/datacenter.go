package server

import (
	datacenterhttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterHttpHandler"
	datacenterPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterPb"
	datacenterrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterRepository"
	datacenterusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterUsecase"
	"fmt"

	grpcconn "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/grpcConn"

	"log"
)

func (s *server) datacenterService() {
	datacenterRepo := datacenterrepository.NewDatacenterRepository(s.db, s.redis)
	datacenterUsecase := datacenterusecase.NewDatacenterUsecase(datacenterRepo, &s.cfg.Grpc)
	datacenterHttpHandler := datacenterhttphandler.NewDatacenterHttpHandler(datacenterUsecase)

	datacenterGrpc := datacenterhttphandler.NewdatacenterGrpcHandler(datacenterUsecase)
	// Grpc client
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.DatacenterUrl)
		datacenterPb.RegisterDataCenterGrpcServiceServer(grpcServer, datacenterGrpc)

		log.Printf("datacenter grpc listening on %s", s.cfg.Grpc.DatacenterUrl)
		grpcServer.Serve(lis)
	}()

	datacenters := s.app.Group("/datacenter_v1")
	datacenters.GET("/get-url", datacenterHttpHandler.FindManyUrlCache)
	datacenters.GET("/get-cache/:cacheId", datacenterHttpHandler.FindCacheData)
	datacenters.POST("/insert-url", datacenterHttpHandler.InsertUrlCache)
	datacenters.DELETE("/delete-url/:urlId", datacenterHttpHandler.DeleteUrlCache)

	s.cron.AddFunc("@every 30s", func() {
		if err := datacenterHttpHandler.CronJobUpdateCache(); err != nil {
			log.Printf("error something wrong %+v", err)
		}
		fmt.Println("Function running every 30 seconds")
	})

}
