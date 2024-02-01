package comment

import "time"

type (
	PushCommentReq struct {
		Title     string    `json:"title" validate:"required"`
		Content   string    `json:"content" validate:"required"`
		CreatedBy string    `json:"created_by"`
		CreateAt  time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	CommentRes struct {
		ProjectId string        `json:"projectId"`
		Comments  []CommentARes `json:"comments"`
		CreateAt  time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
	}

	CommentARes struct {
		CommentId string    `json:"commentId"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		CreatedBy string    `json:"created_by"`
		CreateAt  time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
