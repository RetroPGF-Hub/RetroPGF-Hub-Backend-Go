package project

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProjectModel struct {
		Id           primitive.ObjectID `bson:"_id,omitempty"`
		Name         string             `bson:"name"`
		LogoUrl      string             `bson:"logo_url"`
		GithubUrl    string             `bson:"github_url"`
		WebsiteUrl   string             `bson:"website_url"`
		Description  string             `bson:"description"`
		Feedback     string             `bson:"feedback"`
		Category     string             `bson:"category"`
		FavCount     int64              `bson:"fav_count"`
		CommentCount int64              `bson:"comment_count"`
		CreatedBy    string             `bson:"created_by"`
		CreateAt     time.Time          `bson:"created_at"`
		UpdatedAt    time.Time          `bson:"updated_at"`
	}

	QuestionModel struct {
		Id           primitive.ObjectID `bson:"_id,omitempty"`
		Title        string             `bson:"title"`
		Detail       string             `bson:"detail"`
		CreateAt     time.Time          `bson:"created_at"`
		UpdatedAt    time.Time          `bson:"updated_at"`
		CreatedBy    string             `bson:"created_by"`
		FavCount     int64              `bson:"fav_count"`
		CommentCount int64              `bson:"comment_count"`
	}
)

// name            String
//   logo_url        String?
//   banner_url      String?
//   website_url     String?
//   crypto_category String
//   description     String?   @db.Text
//   reason          String?   @db.Text
//   category        String
//   contact         String?
//   create_by       String
//   user            User      @relation(fields: [create_by], references: [id])
//   Comment         Comment[]
//   create_at       DateTime  @default(now())
//   Like            Like[]
