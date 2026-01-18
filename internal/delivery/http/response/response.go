package response

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"error_code,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, errorCode string, message string, errors interface{}) {
	c.JSON(status, ErrorResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Errors:    errors,
	})
}

func ValidationError(c *gin.Context, err error) {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := toSnakeCase(e.Field())

			switch e.Tag() {
			case "required":
				errors[field] = "Field ini wajib diisi"

			case "min":
				errors[field] = "Field Minimal " + e.Param() + " karakter"

			case "max":
				errors[field] = "Field Maksimal " + e.Param() + " karakter"

			case "len":
				errors[field] = "Field Harus " + e.Param() + " karakter"

			case "email":
				errors[field] = "Format email tidak valid"

			case "oneof":
				errors[field] = "Nilai tidak diperbolehkan"

			case "uuid":
				errors[field] = "Format ID tidak valid"

			case "url":
				errors[field] = "Format URL tidak valid"

			case "numeric":
				errors[field] = "Field Harus berupa angka"

			case "alphanum":
				errors[field] = "Field Hanya boleh huruf dan angka"

			case "gte":
				errors[field] = "Field Harus lebih besar atau sama dengan " + e.Param()

			case "lte":
				errors[field] = "Field Harus lebih kecil atau sama dengan " + e.Param()

			case "eqfield":
				if field == "confirm_password" {
					errors[field] = "Konfirmasi password harus sama dengan password baru"
				} else {
					errors[field] = "Nilai tidak sama dengan field yang dimaksud"
				}

			case "nefield":
				errors[field] = "Nilai tidak boleh sama dengan field lain"

			default:
				errors[field] = "Nilai tidak valid"
			}
		}
	}

	Error(
		c,
		400,
		"VALIDATION_ERROR",
		"Validasi gagal",
		errors,
	)
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
