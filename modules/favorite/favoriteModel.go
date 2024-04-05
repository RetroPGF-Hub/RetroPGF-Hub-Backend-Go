package favorite

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	FavProjectModel struct {
		User      primitive.ObjectID `bson:"_id,omitempty"`
		ProjectId []string           `bson:"projects"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}

	FavQuestionModel struct {
		User       primitive.ObjectID `bson:"_id,omitempty"`
		QuestionId []string           `bson:"questions"`
		CreateAt   time.Time          `bson:"created_at"`
		UpdatedAt  time.Time          `bson:"updated_at"`
	}
)
