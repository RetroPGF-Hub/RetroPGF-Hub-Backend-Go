package datacenterrepository

import "go.mongodb.org/mongo-driver/mongo"

type (
	DatacenterRepositoryService interface {
	}

	datacenterRepository struct {
		db *mongo.Client
	}
)

func NewDatacenterRepository(db *mongo.Client) DatacenterRepositoryService {
	return &datacenterRepository{
		db: db,
	}
}
