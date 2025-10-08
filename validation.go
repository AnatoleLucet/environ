package environ

import (
	"fmt"
	"net/mail"
	"net/url"
	"slices"

	"github.com/AnatoleLucet/as"
)

func validateInt(v string) (int, error) {
	i, err := as.Int(v)
	if err != nil {
		return 0, fmt.Errorf("%w. unable to parse '%s' as integer: %v", ErrInvalidInt, v, err)
	}

	return i, nil
}

func validateFloat(v string) (float64, error) {
	f, err := as.Float(v)
	if err != nil {
		return 0, fmt.Errorf("%w. unable to parse '%s' as float: %v", ErrInvalidFloat, v, err)
	}

	return f, nil
}

func validateBoolean(v string) (bool, error) {
	b, err := as.Bool(v)
	if err != nil {
		return false, fmt.Errorf("%w. unable to parse '%s' as boolean: %v", ErrInvalidBool, v, err)
	}

	return b, nil
}

func validatePort(v string) (int, error) {
	port, err := as.Int(v)
	if err != nil {
		return 0, fmt.Errorf("%w. unable to parse '%s' as integer: %v", ErrInvalidPort, v, err)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("%w. %d is out of range (1-65535)", ErrInvalidPort, port)
	}

	return port, nil
}

func validateUrl(v string) (string, error) {
	if v == "" {
		return "", fmt.Errorf("%w. empty string", ErrInvalidUrl)
	}

	_, err := url.ParseRequestURI(v)
	if err != nil {
		return "", fmt.Errorf("%w. unable to parse '%s' as URL: %v", ErrInvalidUrl, v, err)
	}

	return v, nil
}

func validateEmail(v string) (string, error) {
	if v == "" {
		return "", fmt.Errorf("%w. empty string", ErrInvalidEmail)
	}

	_, err := mail.ParseAddress(v)
	if err != nil {
		return "", fmt.Errorf("%w. unable to parse '%s' as email address: %v", ErrInvalidEmail, v, err)
	}

	return v, nil
}

func validateType[T any](t VariableType, v string) (T, error) {
	var zero T

	switch t {
	case TypeString, "str", "":
		return any(v).(T), nil
	case TypeInt, "integer":
		n, err := validateInt(v)
		return any(n).(T), err
	case TypeFloat:
		f, err := validateFloat(v)
		return any(f).(T), err
	case TypeBoolean, "bool":
		b, err := validateBoolean(v)
		return any(b).(T), err
	case TypePort:
		p, err := validatePort(v)
		return any(p).(T), err
	case TypeUrl:
		u, err := validateUrl(v)
		return any(u).(T), err
	case TypeEmail:
		em, err := validateEmail(v)
		return any(em).(T), err
	}

	return zero, fmt.Errorf("Err: %w. Reason: unknown type '%s'", ErrUnknownType, t)
}

func validate[T comparable](variable Variable[T], value string) (T, error) {
	validated, err := validateType[T](variable.Type, value)
	if err != nil {
		return *new(T), err
	}

	if len(variable.Oneof) > 0 && !slices.Contains(variable.Oneof, validated) {
		return *new(T), fmt.Errorf("%w. Available choices: %v", ErrNotInOneof, variable.Oneof)
	}

	return validated, nil
}
