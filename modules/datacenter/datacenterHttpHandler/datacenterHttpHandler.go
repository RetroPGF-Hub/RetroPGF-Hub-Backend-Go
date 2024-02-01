package datacenterhttphandler

import datacenterusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterUsecase"

type (
	DatacenterHttpHandlerService interface {
	}

	datacenterHttpHandler struct {
		datacenterUsecase datacenterusecase.DatacenterUsecaseService
	}
)

func NewDatacenterHttpHandler(datacenterUsecase datacenterusecase.DatacenterUsecaseService) DatacenterHttpHandlerService {
	return &datacenterHttpHandler{
		datacenterUsecase: datacenterUsecase,
	}
}
