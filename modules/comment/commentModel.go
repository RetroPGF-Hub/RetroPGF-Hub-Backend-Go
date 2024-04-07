package comment

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CommentProjectModel struct {
		ProjectId primitive.ObjectID `bson:"_id,omitempty"`
		Comments  []CommentA         `bson:"comments"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}

	CommentQuestionModel struct {
		QuestionId primitive.ObjectID `bson:"_id,omitempty"`
		Comments   []CommentA         `bson:"comments"`
		CreateAt   time.Time          `bson:"created_at"`
		UpdatedAt  time.Time          `bson:"updated_at"`
	}

	CommentA struct {
		CommentId primitive.ObjectID `bson:"_id,omitempty"`
		Content   string             `bson:"content"`
		CreatedBy primitive.ObjectID `bson:"created_by"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}
)
