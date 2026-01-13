package response

import "github.com/afandimsr/go-gin-api/internal/domain/user"

// Generic success response for swagger
type SuccessUserResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"success"`
	Data    []user.User `json:"data"`
}

type SuccessSingleUserResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"success"`
	Data    user.User `json:"data"`
}

type ErrorSwaggerResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"error"`
	Errors  string `json:"errors"`
}
