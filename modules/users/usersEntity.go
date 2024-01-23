package users

type (
	RegisterUserReq struct {
		Email     string `json:"email" validate:"required,email,max=255"`
		Profile   string `json:"profile" validate:"required,profile"`
		Username  string `json:"username" validate:"required,max=128"`
		Firstname string `json:"firstname" validate:"required,max=128"`
		Lastname  string `json:"lastname" validate:"required,max=128"`
	}

	UserProfileRes struct {
		Id        string `json:"userId"`
		Email     string `json:"email"`
		Profile   string `json:"profile"`
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}
)
