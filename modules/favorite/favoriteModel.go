package favorite

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	FavModel struct {
		User      primitive.ObjectID `bson:"_id,omitempty"`
		ProjectId []string           `bson:"projects"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}
)
