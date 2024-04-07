package project

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"time"
)

type (
	InsertProjectReq struct {
		Name        string `json:"name" validate:"required"`
		Type        string `json:"type" validate:"required"`
		LogoUrl     string `json:"logoUrl" validate:"required"`
		GithubUrl   string `json:"githubUrl" validate:"required"`
		WebsiteUrl  string `json:"websiteUrl" validate:"required"`
		Description string `json:"description" validate:"required"`
		Feedback    string `json:"feedback" validate:"required"`
		Category    string `json:"category" validate:"required"`
		CreatedBy   string `json:"createdBy"`
	}

	InsertQuestionReq struct {
		Name        string `json:"name" validate:"required"`
		Type        string `json:"type" validate:"required"`
		Description string `json:"description" validate:"required"`
		Category    string `json:"category" validate:"required"`
		CreatedBy   string `json:"createdBy"`
	}

	ProjectRes struct {
		Id           string    `json:"_id,omitempty"`
		Type         string    `json:"type"`
		Name         string    `json:"name"`
		LogoUrl      string    `json:"logoUrl"`
		GithubUrl    string    `json:"githubUrl"`
		WebsiteUrl   string    `json:"websiteUrl"`
		Description  string    `json:"description"`
		Feedback     string    `json:"feedback"`
		Category     string    `json:"category"`
		FavCount     int64     `json:"favCount"`
		CommentCount int64     `json:"commentCount"`
		CreatedBy    string    `json:"createdBy"`
		FavOrNot     bool      `json:"favOrNot"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	ProjectResWithUser struct {
		Id           string                  `json:"_id,omitempty"`
		Type         string                  `json:"type"`
		Name         string                  `json:"name"`
		LogoUrl      string                  `json:"logoUrl"`
		GithubUrl    string                  `json:"githubUrl"`
		WebsiteUrl   string                  `json:"websiteUrl"`
		Description  string                  `json:"description"`
		Feedback     string                  `json:"feedback"`
		Category     string                  `json:"category"`
		FavCount     int64                   `json:"favCount"`
		CommentCount int64                   `json:"commentCount"`
		FavOrNot     bool                    `json:"favOrNot"`
		Owner        users.SecureUserProfile `json:"owner"`
		CreatedAt    time.Time               `json:"createdAt"`
		UpdatedAt    time.Time               `json:"updatedAt"`
	}

	FullProjectRes struct {
		Id           string                        `json:"_id,omitempty"`
		Name         string                        `json:"name"`
		Type         string                        `json:"type"`
		LogoUrl      string                        `json:"logoUrl"`
		GithubUrl    string                        `json:"githubUrl"`
		WebsiteUrl   string                        `json:"websiteUrl"`
		Description  string                        `json:"description"`
		Feedback     string                        `json:"feedback"`
		Category     string                        `json:"category"`
		FavCount     int64                         `json:"favCount"`
		CommentCount int64                         `json:"commentCount"`
		FavOrNot     bool                          `json:"favOrNot"`
		Owner        users.SecureUserProfile       `json:"owner"`
		Comment      []comment.CommentAResWithUser `json:"comment"`
		CreateAt     time.Time                     `json:"createdAt"`
		UpdatedAt    time.Time                     `json:"updatedAt"`
	}

	RandomProjectDisplay struct {
		Id           string `json:"_id,omitempty"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		LogoUrl      string `json:"logoUrl"`
		Category     string `json:"category"`
		Description  string `json:"description"`
		FavCount     int64  `json:"favCount"`
		CommentCount int64  `json:"commentCount"`
	}
)
