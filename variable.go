package environ

import (
	"fmt"
	"os"
)

type VariableType string

var (
	TypeString  VariableType = "string"
	TypeInt     VariableType = "int"
	TypeFloat   VariableType = "float"
	TypeBoolean VariableType = "boolean"
	TypePort    VariableType = "port"
	TypeUrl     VariableType = "url"
	TypeEmail   VariableType = "email"
)

type VariableValidator[T comparable] func(T) (T, error)

type Variable[T comparable] struct {
	Name        string       `tag:"env | get('name')"`
	Type        VariableType `tag:"env | get('type')"`
	Default     *T           `tag:"env | get('default')"`
	Optional    bool         `tag:"env | has('optional')"`
	Description string       `tag:"env | get('desc')"`
	Oneof       []T          `tag:"env | get('oneof') | split('|')"`

	Validator VariableValidator[T] `env:"-"`
}

// Load will fetch the environment variable, validate it, and return the value or an error
func (v Variable[T]) Load() (T, error) {
	validated, err := loadVariable(v)
	if err != nil {
		return *new(T), fmt.Errorf("Err: variable %q. Reason: %w", v.Name, err)
	}

	return validated, nil
}

// MustLoad is like Load but will panic if there is an error
func (v Variable[T]) MustLoad() T {
	validated, err := v.Load()
	if err != nil {
		panic(err)
	}

	return validated
}

type VariableBuilder[T comparable] struct {
	Variable[T]
}

func (vb VariableBuilder[T]) Optional() VariableBuilder[T] {
	vb.Variable.Optional = true
	return vb
}

func (vb VariableBuilder[T]) Oneof(choices ...T) VariableBuilder[T] {
	vb.Variable.Oneof = choices
	return vb
}

func (vb VariableBuilder[T]) Default(value T) VariableBuilder[T] {
	vb.Variable.Default = &value
	return vb
}

func (vb VariableBuilder[T]) Desc(desc string) VariableBuilder[T] {
	vb.Variable.Description = desc
	return vb
}

func (vb VariableBuilder[T]) Validate(validator VariableValidator[T]) VariableBuilder[T] {
	vb.Variable.Validator = validator
	return vb
}

func loadVariable[T comparable](variable Variable[T]) (T, error) {
	if variable.Name == "" {
		return *new(T), ErrMissingName
	}

	value, exists := os.LookupEnv(variable.Name)
	if !exists || value == "" {
		if variable.Default != nil {
			return *variable.Default, nil
		} else if variable.Optional {
			return *new(T), nil
		} else {
			return *new(T), ErrMissingValue
		}
	}

	validated, err := validate(variable, value)
	if err != nil {
		return *new(T), err
	}

	if variable.Validator != nil {
		validated, err = variable.Validator(validated)
		if err != nil {
			return *new(T), err
		}
	}

	return validated, nil
}
