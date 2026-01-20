package middleware

import (
	"errors"
	"net/http"

	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/response"
	"github.com/afandimsr/go-gin-api/internal/domain/apperror"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		isDevelopment := config.AppEnv == "development" || gin.Mode() == gin.DebugMode

		var appErr *apperror.AppError
		var details interface{} = nil

		if errors.As(err, &appErr) {

			// 1️⃣ VALIDATION DARI GIN (binding / request)
			if appErr.ErrorCode == apperror.ValidationError {
				if appErr.Err != nil {
					// validator.ValidationErrors
					response.ValidationError(c, appErr.Err)
					c.Abort()
					return
				}
			}

			// 2️⃣ DOMAIN VALIDATION (business rule)
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
				c.Abort()
				return
			}

			// 3️⃣ BUSINESS ERROR (NotFound, Unauthorized, etc)
			var messages string

			messages = appErr.Message
			if !isDevelopment {
				messages = response.MessagesMap[apperror.GeneralError]
			}

			response.Error(
				c,
				appErr.Code,
				appErr.ErrorCode,
				messages,
				nil,
			)
			c.Abort()
			return
		}

		// 4️⃣ FALLBACK SYSTEM ERROR
		message := "internal server error"

		if isDevelopment {
			message = err.Error()
			details = map[string]interface{}{
				"error": err.Error(),
				"type":  errors.Unwrap(err),
			}
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			apperror.SystemInternalError,
			message,
			details,
		)
		c.Abort()
	}
}
