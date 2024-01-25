package project

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProjectModel struct {
		Id             primitive.ObjectID `bson:"_id,omitempty"`
		Name           string             `bson:"name"`
		LogoUrl        string             `bson:"logo_url"`
		BannerUrl      string             `bson:"banner_url"`
		WebsiteUrl     string             `bson:"website_url"`
		CryptoCategory string             `bson:"crypto_category"`
		Description    string             `bson:"description"`
		Reason         string             `bson:"reason"`
		Category       string             `bson:"category"`
		Contact        string             `bson:"contact"`
		CreatedBy      primitive.ObjectID `bson:"created_by"`
		CreateAt       time.Time          `bson:"created_at"`
		UpdatedAt      time.Time          `bson:"updated_at"`
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
