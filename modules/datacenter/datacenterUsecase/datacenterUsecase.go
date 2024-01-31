package datacenterusecase

import datacenterrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterRepository"

type (
	DatacenterUsecaseService interface {
	}

	datacenterUsecase struct {
		datacenterRepo datacenterrepository.DatacenterRepositoryService
	}
)

func NewDatacenterUsecase(datacenterRepo datacenterrepository.DatacenterRepositoryService) DatacenterUsecaseService {
	return &datacenterUsecase{
		datacenterRepo: datacenterRepo,
	}
}
