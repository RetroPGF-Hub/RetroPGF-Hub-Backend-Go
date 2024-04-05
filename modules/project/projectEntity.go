package project

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"time"
)

type (
	InsertProjectReq struct {
		Name        string `json:"name" validate:"required"`
		LogoUrl     string `json:"logoUrl" validate:"required"`
		GithubUrl   string `json:"githubUrl" validate:"required"`
		WebsiteUrl  string `json:"websiteUrl" validate:"required"`
		Description string `json:"description" validate:"required"`
		Feedback    string `json:"feedback" validate:"required"`
		Category    string `json:"category" validate:"required"`
		CreatedBy   string `json:"createdBy"`
	}

	ProjectRes struct {
		Id           string    `json:"_id,omitempty"`
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
		Id           string               `json:"_id,omitempty"`
		Name         string               `json:"name"`
		LogoUrl      string               `json:"logoUrl"`
		GithubUrl    string               `json:"githubUrl"`
		WebsiteUrl   string               `json:"websiteUrl"`
		Description  string               `json:"description"`
		Feedback     string               `json:"feedback"`
		Category     string               `json:"category"`
		FavCount     int64                `json:"favCount"`
		CommentCount int64                `json:"commentCount"`
		FavOrNot     bool                 `json:"favOrNot"`
		Owner        users.UserProfileRes `json:"owner"`
		CreatedAt    time.Time            `json:"createdAt"`
		UpdatedAt    time.Time            `json:"updatedAt"`
	}

	FullProjectRes struct {
		Id           string                        `json:"_id,omitempty"`
		Name         string                        `json:"name"`
		LogoUrl      string                        `json:"logoUrl"`
		GithubUrl    string                        `json:"githubUrl"`
		WebsiteUrl   string                        `json:"websiteUrl"`
		Description  string                        `json:"description"`
		Feedback     string                        `json:"feedback"`
		Category     string                        `json:"category"`
		FavCount     int64                         `json:"favCount"`
		CommentCount int64                         `json:"commentCount"`
		FavOrNot     bool                          `json:"favOrNot"`
		Owner        users.UserProfileRes          `json:"owner"`
		Comment      []comment.CommentAResWithUser `json:"comment"`
		CreateAt     time.Time                     `json:"createdAt"`
		UpdatedAt    time.Time                     `json:"updatedAt"`
	}

	InsertQuestionReq struct {
		Title     string `json:"title" validate:"required"`
		Detail    string `json:"detail" validate:"required"`
		CreatedBy string `json:"created_by" validate:"required"`
	}

	QuestionRes struct {
		Id           string    `json:"_id,omitempty"`
		Title        string    `json:"title"`
		Detail       string    `json:"detail"`
		CreateAt     time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		CreatedBy    string    `json:"createdBy"`
		FavCount     int64     `json:"fav_count"`
		CommentCount int64     `json:"comment_count"`
	}
	QuestionResWithUser struct {
		Id           string               `json:"_id,omitempty"`
		Title        string               `json:"title"`
		Detail       string               `json:"detail"`
		CreateAt     time.Time            `json:"created_at"`
		UpdatedAt    time.Time            `json:"updated_at"`
		Owner        users.UserProfileRes `json:"owner"`
		FavCount     int64                `json:"fav_count"`
		CommentCount int64                `json:"comment_count"`
	}

	FullQuestionRes struct {
		Id           string                        `json:"_id,omitempty"`
		Title        string                        `json:"title"`
		Detail       string                        `json:"detail"`
		CreateAt     time.Time                     `json:"created_at"`
		UpdatedAt    time.Time                     `json:"updated_at"`
		Owner        users.UserProfileRes          `json:"owner"`
		Comment      []comment.CommentAResWithUser `json:"comment"`
		FavCount     int64                         `json:"fav_count"`
		CommentCount int64                         `json:"comment_count"`
	}
)
