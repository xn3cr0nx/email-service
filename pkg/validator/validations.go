package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// TestingValidator to test custom validators
func TestingValidator(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) > 5
}

// LimitValidator to test custom validators
func LimitValidator(fl validator.FieldLevel) bool {
	return fl.Field().Int() < 500 && fl.Field().Int()%5 == 0
}

// PasswordValidator to test custom validators
func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	includesNumber, _ := regexp.MatchString(".*[0-9].*", password)
	includesUpper, _ := regexp.MatchString(".*[A-Z].*", password)
	return len(password) > 6 && includesNumber && includesUpper
}
