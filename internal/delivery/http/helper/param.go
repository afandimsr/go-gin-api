package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/afandimsr/go-gin-api/internal/domain/apperror"
)

func ValidateUUIDParam(c *gin.Context, name string) (string, error) {
	val := c.Param(name)
	if _, err := uuid.Parse(val); err != nil {
		return "", apperror.Validation(err).
			WithCode(apperror.ValidationInvalidUUID)
	}
	return val, nil
}

func ValidateUUID(val string) error {
	if _, err := uuid.Parse(val); err != nil {
		return apperror.Validation(err).
			WithCode(apperror.ValidationInvalidUUID)
	}
	return nil
}

func ValidateUUIDParamNotFound(c *gin.Context, name string) (string, error) {
	val := c.Param(name)

	if _, err := uuid.Parse(val); err != nil {
		return "", apperror.NotFound(
			"resource not found",
			err,
		).WithCode(apperror.ResourceNotFound)
	}

	return val, nil
}
