package favorite

type (
	FavReq struct {
		ProjectId string `json:"projectId" validate:"required"`
		User      string `json:"user" validate:"required"`
	}
)
