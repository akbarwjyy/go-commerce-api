package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	cv := New()
	v := cv.GetValidator()

	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"Valid password", "Password123", true},
		{"Too short", "Pass1", false},
		{"No uppercase", "password123", false},
		{"No lowercase", "PASSWORD123", false},
		{"No number", "Passwordabc", false},
		{"Empty", "", false},
	}

	type TestStruct struct {
		Password string `validate:"password"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(&TestStruct{Password: tt.password})
			if tt.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	cv := New()
	v := cv.GetValidator()

	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{"Valid +62", "+6281234567890", true},
		{"Valid 62", "6281234567890", true},
		{"Valid 08", "081234567890", true},
		{"Too short", "0812345", false},
		{"Invalid format", "12345678901", false},
		{"Empty", "", false},
	}

	type TestStruct struct {
		Phone string `validate:"phone"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(&TestStruct{Phone: tt.phone})
			if tt.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestValidateNoSpaces(t *testing.T) {
	cv := New()
	v := cv.GetValidator()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"No spaces", "username", true},
		{"With spaces", "user name", false},
		{"Empty", "", true},
	}

	type TestStruct struct {
		Username string `validate:"no_spaces"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(&TestStruct{Username: tt.value})
			if tt.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestValidateAlphaSpace(t *testing.T) {
	cv := New()
	v := cv.GetValidator()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Letters only", "JohnDoe", true},
		{"Letters with space", "John Doe", true},
		{"With numbers", "John123", false},
		{"With special chars", "John@Doe", false},
	}

	type TestStruct struct {
		Name string `validate:"alpha_space"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(&TestStruct{Name: tt.value})
			if tt.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestGetErrorMessage(t *testing.T) {
	msg := ValidationErrorMessages["required"]
	assert.Equal(t, "This field is required", msg)

	msg = ValidationErrorMessages["email"]
	assert.Equal(t, "Invalid email format", msg)

	msg = ValidationErrorMessages["password"]
	assert.Equal(t, "Password must be at least 8 characters with uppercase, lowercase, and number", msg)
}
