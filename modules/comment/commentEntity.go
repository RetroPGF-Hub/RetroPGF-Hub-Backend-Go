package comment

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"time"
)

type (
	PushCommentReq struct {
		Content   string    `json:"content" validate:"required"`
		CreatedBy string    `json:"created_by"`
		CreateAt  time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	CommentProjectRes struct {
		ProjectId string        `json:"projectId"`
		Comments  []CommentARes `json:"comments"`
		CreateAt  time.Time     `json:"createdAt"`
		UpdatedAt time.Time     `json:"updatedAt"`
	}

	CommentQuestionRes struct {
		QuestionId string        `json:"questionId"`
		Comments   []CommentARes `json:"comments"`
		CreateAt   time.Time     `json:"createdAt"`
		UpdatedAt  time.Time     `json:"updatedAt"`
	}

	CommentARes struct {
		CommentId string    `json:"commentId"`
		Content   string    `json:"content"`
		CreatedBy string    `json:"createdBy"`
		CreateAt  time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	CommentAResWithUser struct {
		CommentId string                  `json:"commentId"`
		Content   string                  `json:"content"`
		CreatedBy users.SecureUserProfile `json:"createdBy"`
		CreateAt  time.Time               `json:"createdAt"`
		UpdatedAt time.Time               `json:"updatedAt"`
	}
)
