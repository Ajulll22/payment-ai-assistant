package validation

import (
	"strings"

	"github.com/Ajulll22/payment-ai-assistant/pkg/formatter"
	"github.com/go-playground/validator/v10"
)

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "required_if":
		return "This field is required if " + fe.Param()
	case "required_without":
		return "This field is required if " + strings.Replace(fe.Param(), " ", " or ", -1) + " is not present"
	case "lt":
		return "Should be less than " + fe.Param()
	case "lte":
		return "Should be less than or equal to " + fe.Param()
	case "gt":
		return "Should be greater than " + fe.Param()
	case "gte":
		return "Should be greater than or equal to " + fe.Param()
	case "ltfield":
		return "Should be less than " + formatter.ToTitleCase(fe.Param())
	case "ltefield":
		return "Should be less than or equal " + formatter.ToTitleCase(fe.Param())
	case "gtfield":
		return "Should be greater than " + formatter.ToTitleCase(fe.Param())
	case "gtefield":
		return "Should be greater than or equal " + formatter.ToTitleCase(fe.Param())
	case "nefield":
		return "Shoult not equal to " + formatter.ToTitleCase(fe.Param())
	case "numeric":
		return "This field should be numeric"
	case "email":
		return "This field should be a valid email address"
	case "max":
		return "Length maximum " + fe.Param()
	case "min":
		return "Length minimum " + fe.Param()
	case "oneof":
		return "Should be one of " + strings.Replace(fe.Param(), " ", ", ", -1)
	case "unique":
		return "Field " + formatter.ToTitleCase(fe.Param()) + " can't have the same value"
	case "filesize":
		return "File size must be less than " + formatter.ToSnakeCase(fe.Param()) + "MB"
	case "filetype":
		return "File type does not match with " + formatter.ToSnakeCase(fe.Param()) + " type"
	}

	return "Unknown error"

}

func FormatValidation(ve validator.ValidationErrors) []string {
	errList := []string{}

	for _, fe := range ve {
		errList = append(errList, formatter.ToTitleCase(fe.Field())+", "+getErrorMsg(fe))
	}

	return errList
}
