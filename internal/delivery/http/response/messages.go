package response

type FieldError struct {
	Field   string
	Message string
}

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
