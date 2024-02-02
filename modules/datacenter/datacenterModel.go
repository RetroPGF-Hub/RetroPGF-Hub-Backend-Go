package datacenter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CacheModel struct {
		UrlId     primitive.ObjectID `bson:"_id,omitempty" json:"urlId"`
		Url       string             `bson:"url" json:"url" validate:"required"`
		CreateAt  time.Time          `bson:"created_at" json:"createdAt"`
		UpdatedAt time.Time          `bson:"updated_at" json:"ureatedAt"`
	}
)
