package middleware

import (
	"errors"
	"net/http"

	"github.com/afandimsr/go-gin-api/internal/delivery/http/response"
	"github.com/afandimsr/go-gin-api/internal/domain/apperror"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {

			// 1️⃣ VALIDATION DARI GIN (binding / request)
			if appErr.ErrorCode == apperror.ValidationError {
				if appErr.Err != nil {
					// validator.ValidationErrors
					response.ValidationError(c, appErr.Err)
					return
				}
			}

			// 2️⃣ DOMAIN VALIDATION (Password, business rule)
			if fieldErr, ok := response.ValidationErrorMap[appErr.ErrorCode]; ok {
				response.Error(
					c,
					appErr.Code,
					appErr.ErrorCode,
					appErr.Message,
					map[string]string{
						fieldErr.Field: fieldErr.Message,
					},
				)
				return
			}

			// 3️⃣ BUSINESS ERROR (NotFound, Unauthorized, etc)
			response.Error(
				c,
				appErr.Code,
				appErr.ErrorCode,
				appErr.Message,
				nil,
			)
			return
		}

		// 4️⃣ FALLBACK SYSTEM ERROR
		response.Error(
			c,
			http.StatusInternalServerError,
			apperror.SystemInternalError,
			"internal server error",
			nil,
		)
	}
}
