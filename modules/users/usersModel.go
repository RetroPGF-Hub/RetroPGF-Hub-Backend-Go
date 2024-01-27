package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserDb struct {
		Id        primitive.ObjectID `bson:"_id,omitempty"`
		Source    string             `bson:"source"`
		Email     string             `bson:"email"`
		Profile   string             `bson:"profile"`
		Password  string             `bson:"password"`
		Username  string             `bson:"user_name"`
		Firstname string             `bson:"first_name"`
		Lastname  string             `bson:"last_name"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}

	Credential struct {
		Id           primitive.ObjectID `bson:"_id,omitempty"`
		UserId       string             `bson:"user_id"`
		AccessToken  string             `bson:"access_token"`
		RefreshToken string             `bson:"refresh_token"`
		CreatedAt    time.Time          `bson:"created_at"`
		UpdatedAt    time.Time          `bson:"updated_at"`
	}

	UpdateRefreshTokenReq struct {
		UserId       string    `bson:"user_id"`
		AccessToken  string    `bson:"access_token"`
		RefreshToken string    `bson:"refresh_token"`
		UpdatedAt    time.Time `bson:"updated_at"`
	}
)
