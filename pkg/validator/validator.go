package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator with custom validations
type CustomValidator struct {
	validate *validator.Validate
}

// New creates a new CustomValidator with custom validations registered
func New() *CustomValidator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("no_spaces", validateNoSpaces)
	v.RegisterValidation("alpha_space", validateAlphaSpace)

	return &CustomValidator{validate: v}
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

// GetValidator returns the underlying validator
func (cv *CustomValidator) GetValidator() *validator.Validate {
	return cv.validate
}

// validatePassword checks password strength:
// - Minimum 8 characters
// - At least 1 uppercase letter
// - At least 1 lowercase letter
// - At least 1 number
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// validatePhone checks phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Indonesian phone format: +62xxx or 08xxx
	phoneRegex := regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,10}$`)
	return phoneRegex.MatchString(phone)
}

// validateNoSpaces ensures no spaces in the string
func validateNoSpaces(fl validator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), " ")
}

// validateAlphaSpace allows only alphabets and spaces
func validateAlphaSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	alphaSpaceRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return alphaSpaceRegex.MatchString(value)
}

// ValidationErrorMessages provides custom error messages
var ValidationErrorMessages = map[string]string{
	"required":    "This field is required",
	"email":       "Invalid email format",
	"min":         "Value is too short",
	"max":         "Value is too long",
	"gte":         "Value must be greater than or equal to the minimum",
	"lte":         "Value must be less than or equal to the maximum",
	"gt":          "Value must be greater than zero",
	"password":    "Password must be at least 8 characters with uppercase, lowercase, and number",
	"phone":       "Invalid phone number format",
	"no_spaces":   "Spaces are not allowed",
	"alpha_space": "Only alphabets and spaces are allowed",
	"oneof":       "Invalid value",
	"url":         "Invalid URL format",
}

// GetErrorMessage returns a human-readable error message
func GetErrorMessage(fe validator.FieldError) string {
	if msg, ok := ValidationErrorMessages[fe.Tag()]; ok {
		return msg
	}
	return fe.Error()
}

// FormatValidationErrors formats validation errors to a map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			field := strings.ToLower(fe.Field())
			errors[field] = GetErrorMessage(fe)
		}
	}

	return errors
}
