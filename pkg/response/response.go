package response

import "github.com/labstack/echo/v4"

type (
	MsgResponse struct {
		Message string `json:"msg"`
	}
)

func ErrResponse(c echo.Context, status int, msg string) error {
	return c.JSON(status, &MsgResponse{Message: msg})
}

func SuccessResponse(c echo.Context, status int, data any) error {
	return c.JSON(status, data)
}
