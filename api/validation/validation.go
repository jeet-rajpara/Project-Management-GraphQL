package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)


func ValidateInput(input interface{}) error {
	validate := validator.New()

	// Register any custom validation tags or functions here if needed
	validate.RegisterValidation("notblank", IsNotJustWhitespace)

	// Perform validation
	if err := validate.Struct(input); err != nil {
		// Validation failed
		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' validation failed on tag '%s'", err.Field(), err.Tag()))
		}
		return fmt.Errorf(strings.Join(errMsgs, ", "))
	}

	return nil
}

func IsNotJustWhitespace(fl validator.FieldLevel) bool {
    str := fl.Field().String()
    return len(strings.TrimSpace(str)) > 0
}
