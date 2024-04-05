package favorite

type (
	FavProjectReq struct {
		ProjectId string `json:"projectId" validate:"required"`
		User      string `json:"user" validate:"required"`
	}
	FavQuestionReq struct {
		QuestionId string `json:"questionid" validate:"required"`
		User       string `json:"user" validate:"required"`
	}
)
