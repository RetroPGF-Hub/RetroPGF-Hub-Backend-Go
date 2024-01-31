package server

import (
	datacenterhttphandler "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterHttpHandler"
	datacenterrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterRepository"
	datacenterusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterUsecase"
)

func (s *server) datacenterService() {
	datacenterRepo := datacenterrepository.NewDatacenterRepository(s.db)
	datacenterUsecase := datacenterusecase.NewDatacenterUsecase(datacenterRepo)
	datacenterHttpHandler := datacenterhttphandler.NewDatacenterHttpHandler(datacenterUsecase)

	_ = datacenterHttpHandler
}
