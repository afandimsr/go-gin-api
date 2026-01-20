package response

type FieldError struct {
	Field   string
	Message string
}

// ValidationErrorMap maps specific validation error codes to their corresponding field errors
var ValidationErrorMap = map[string]FieldError{
	"PASSWORD_TOO_SHORT": {
		Field:   "new_password",
		Message: "Password minimal 8 karakter",
	},
	"PASSWORD_WEAK": {
		Field:   "new_password",
		Message: "Password harus mengandung huruf besar, huruf kecil, angka, dan simbol",
	},
	"VALIDATION_INVALID_UUID": {
		Field:   "id",
		Message: "Format ID tidak valid",
	},
}

type MessageResponse struct {
	Message string
}

// MessagesMap contains predefined messages for various error codes and success messages
var MessagesMap = map[string]string{
	"GENERAL_ERROR": "Terjadi kesalahan pada sistem, silakan coba beberapa saat lagi", // required in production
}
