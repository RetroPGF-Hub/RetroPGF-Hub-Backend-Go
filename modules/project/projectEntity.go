package project

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"time"
)

type (
	InsertProjectReq struct {
		Name           string `json:"name" validate:"required"`
		LogoUrl        string `json:"logoUrl" validate:"required"`
		BannerUrl      string `json:"bannerUrl" validate:"required"`
		WebsiteUrl     string `json:"websiteUrl" validate:"required"`
		CryptoCategory string `json:"cryptoCategory" validate:"required"`
		Description    string `json:"description" validate:"required"`
		Reason         string `json:"reason" validate:"required"`
		Category       string `json:"category" validate:"required"`
		Contact        string `json:"contact" validate:"required"`
		CreatedBy      string `json:"createdBy"`
	}

	ProjectRes struct {
		Id             string    `json:"_id,omitempty"`
		Name           string    `json:"name"`
		LogoUrl        string    `json:"logoUrl"`
		BannerUrl      string    `json:"bannerUrl"`
		WebsiteUrl     string    `json:"websiteUrl"`
		CryptoCategory string    `json:"cryptoCategory"`
		Description    string    `json:"description"`
		Reason         string    `json:"reason"`
		Category       string    `json:"category"`
		Contact        string    `json:"contact"`
		FavCount       int64     `json:"favCount"`
		CommentCount   int64     `json:"commentCount"`
		CreatedBy      string    `json:"createdBy"`
		FavOrNot       bool      `json:"favOrNot"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}

	ProjectResWithUser struct {
		Id             string               `json:"_id,omitempty"`
		Name           string               `json:"name"`
		LogoUrl        string               `json:"logoUrl"`
		BannerUrl      string               `json:"bannerUrl"`
		WebsiteUrl     string               `json:"websiteUrl"`
		CryptoCategory string               `json:"cryptoCategory"`
		Description    string               `json:"description"`
		Reason         string               `json:"reason"`
		Category       string               `json:"category"`
		Contact        string               `json:"contact"`
		FavCount       int64                `json:"favCount"`
		CommentCount   int64                `json:"commentCount"`
		FavOrNot       bool                 `json:"favOrNot"`
		Owner          users.UserProfileRes `json:"owner"`
		CreatedAt      time.Time            `json:"createdAt"`
		UpdatedAt      time.Time            `json:"updatedAt"`
	}

	FullProjectRes struct {
		Id             string                        `json:"_id,omitempty"`
		Name           string                        `json:"name"`
		LogoUrl        string                        `json:"logoUrl"`
		BannerUrl      string                        `json:"bannerUrl"`
		WebsiteUrl     string                        `json:"websiteUrl"`
		CryptoCategory string                        `json:"cryptoCategory"`
		Description    string                        `json:"description"`
		Reason         string                        `json:"reason"`
		Category       string                        `json:"category"`
		Contact        string                        `json:"contact"`
		FavCount       int64                         `json:"favCount"`
		CommentCount   int64                         `json:"commentCount"`
		FavOrNot       bool                          `json:"favOrNot"`
		Owner          users.UserProfileRes          `json:"owner"`
		Comment        []comment.CommentAResWithUser `json:"comment"`
		CreateAt       time.Time                     `json:"createdAt"`
		UpdatedAt      time.Time                     `json:"updatedAt"`
	}
)
