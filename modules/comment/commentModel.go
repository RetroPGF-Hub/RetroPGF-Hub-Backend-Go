package comment

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CommentModel struct {
		ProjectId primitive.ObjectID `bson:"_id,omitempty"`
		Comments  []CommentA         `bson:"comments"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}

	CommentA struct {
		CommentId primitive.ObjectID `bson:"_id,omitempty"`
		Title     string             `bson:"title"`
		Content   string             `bson:"content"`
		CreatedBy string             `bson:"created_by"`
		CreateAt  time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`
	}
)
