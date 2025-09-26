package validation

import (
	"github.com/go-playground/validator/v10"
)

var Val *validator.Validate

func init() {
	Val = validator.New()
}
