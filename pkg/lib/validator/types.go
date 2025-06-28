package validator

import (
	"errors"
	"fmt"
)

type ServiceError struct {
	Code          string
	OriginMessage string
	Message       string
	Attributes    map[string]string
}

func (e ServiceError) Error() string {
	return e.Message
}

func (e ServiceError) Is(tgt error) bool {
	target := ServiceError{}
	ok := errors.As(tgt, &target)
	if !ok {
		return false
	}

	return e.Message == target.Message && e.Code == target.Code
}

func SvcError(errInt int, metadata map[string]string, customMessage string) ServiceError {
	code := "9000000"
	message := ""
	lemonCode, ok := ErrorCodeMap[errInt]
	if ok {
		code = lemonCode
	}

	if customMessage != "" {
		message = customMessage
	}

	return ServiceError{
		Code:       code,
		Message:    message,
		Attributes: metadata,
	}
}

const (
	ErrGeneralError int = iota
	ErrInvalidMandatoryField
	ErrInvalidFormatField
)

const (
	ErrDomainGeneralPrefix = "9999"
)

var (
	ErrorCodeMap = map[int]string{
		// domain + code category + incr

		// general
		ErrGeneralError:          fmt.Sprintf("%s%s", ErrDomainGeneralPrefix, "5000"),
		ErrInvalidMandatoryField: fmt.Sprintf("%s%s", ErrDomainGeneralPrefix, "4001"),
		ErrInvalidFormatField:    fmt.Sprintf("%s%s", ErrDomainGeneralPrefix, "4002"),
	}
)
