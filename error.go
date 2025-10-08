package environ

import "errors"

var (
	ErrInvalidPort  = errors.New("invalid port")
	ErrInvalidUrl   = errors.New("invalid url")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidBool  = errors.New("invalid boolean value")
	ErrInvalidInt   = errors.New("invalid int")
	ErrInvalidFloat = errors.New("invalid float")
	ErrUnknownType  = errors.New("unknown variable type")

	ErrNotInOneof   = errors.New("the value is not a possible choice")
	ErrMissingValue = errors.New("missing required variable")

	ErrMissingName     = errors.New("missing variable name")
	ErrInvalidTag      = errors.New("invalid variable tag")
	ErrSetField        = errors.New("field is not settable")
	ErrUnsupportedType = errors.New("unsupported variable type")

	ErrUnexpected = errors.New("unexpected error")
)
