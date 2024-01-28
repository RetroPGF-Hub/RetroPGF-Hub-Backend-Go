package favorite

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	FavModel struct {
		ProjectId primitive.ObjectID `bson:"_id,omitempty"`
		Users     []string           `bson:"users"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}
)
