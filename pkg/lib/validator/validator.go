package validator

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/fatih/structs"
	"strings"
	"unicode"
)

func ValidateStruct(s interface{}) error {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		validatorErr, isRequired := ExtractGoValidatorError(s, err)

		message := fmt.Sprintf("error invalid field: %s", validatorErr.Error())
		if isRequired {
			message = fmt.Sprintf("error invalid mandatory field: %s", validatorErr.Name)
			return SvcError(
				ErrInvalidMandatoryField,
				map[string]string{"field_name": validatorErr.Error()},
				message,
			)
		}

		return SvcError(
			ErrInvalidFormatField,
			map[string]string{"field_name": validatorErr.Error()},
			message,
		)
	}
	return nil
}

func ExtractGoValidatorError(obj interface{}, err error) (valErr govalidator.Error, isRequired bool) {
	if err == nil {
		return valErr, isRequired
	}

	errors := err.(govalidator.Errors)
	valErr, ok := errors.Errors()[0].(govalidator.Error)
	if !ok {
		return ExtractGoValidatorError(obj, errors.Errors()[0])
	}

	var field *structs.Field
	s := structs.New(obj)

	if len(valErr.Path) == 0 {
		// validatorErr.Name first char is lower case
		fieldName := []rune(valErr.Name)
		fieldName[0] = unicode.ToUpper(fieldName[0])
		field, _ = s.FieldOk(string(fieldName))
	} else {
		// currently only support one nested struct field
		fieldName := []rune(valErr.Name)
		fieldName[0] = unicode.ToUpper(fieldName[0])
		field = s.Field(valErr.Path[0]).Field(string(fieldName))
	}

	validatorTags := field.Tag("valid")
	isRequired = strings.Contains(validatorTags, "required")

	return valErr, isRequired

}
