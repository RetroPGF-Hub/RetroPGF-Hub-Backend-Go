package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserDb struct {
		Id        primitive.ObjectID `bson:"_id,omitempty"`
		Email     string             `bson:"email"`
		Profile   string             `bson:"profile"`
		Username  string             `bson:"user_name"`
		Firstname string             `bson:"first_name"`
		Lastname  string             `bson:"last_name"`
		CreateAt  time.Time          `bson:"create_at"`
		UpdatedAt time.Time          `bson:"create_at"`
	}
)
