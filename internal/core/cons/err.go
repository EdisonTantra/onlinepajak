package cons

import "errors"

var (
	ErrInvalidNameLength     = errors.New("error invalid name length")
	ErrInvalidPhoneLength    = errors.New("error invalid phone length")
	ErrInvalidPhonePrefix    = errors.New("error invalid phone prefix")
	ErrInvalidPasswordLength = errors.New("error invalid password length")
	ErrInvalidPasswordFormat = errors.New("error password must have capital, number and symbol")
	ErrInvalidAuthorized     = errors.New("user not authorize for this action")
	ErrInvalidToken          = errors.New("error invalid token")
	ErrDataConflict          = errors.New("error data conflict")

	ErrJWTFormat     = errors.New("invalid token format")
	ErrJWTSign       = errors.New("invalid signature")
	ErrJWTSignMethod = errors.New("invalid signing method")
	ErrJWTExpired    = errors.New("invalid token expired")
)

const (
	MinLengthName  = 3
	MaxLengthName  = 60
	MinLengthPass  = 6
	MaxLengthPass  = 64
	MinLengthPhone = 10
	MaxLengthPhone = 13
	PrefixPhoneID  = "+62"
	AuthTokenType  = "Bearer"
)
