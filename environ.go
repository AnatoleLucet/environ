package environ

import (
	"errors"
	"fmt"

	"github.com/AnatoleLucet/tiq"
)

func Load[T any]() (T, error) {
	return load[T]()
}

func MustLoad[T any]() T {
	t, err := load[T]()
	if err != nil {
		panic(err)
	}

	return t
}

func load[T any]() (T, error) {
	var t T

	inspector, err := tiq.Inspect(&t)
	if err != nil {
		return t, fmt.Errorf("%w: %v", ErrUnsupportedType, err)
	}

	for _, field := range inspector.Fields() {
		variable, err := tiq.Parse[Variable[any]](field)
		if err != nil {
			if errors.Is(err, tiq.ErrCompileTag) {
				return t, fmt.Errorf("%w for field %q: %v", ErrInvalidTag, field.Name, err)
			}

			return t, fmt.Errorf("%w for field %q: %v", ErrUnexpected, field.Name, err)
		}

		if variable.Name == "" {
			continue
		}

		value, err := variable.Load()
		if err != nil {
			return t, err
		}

		if value == nil {
			continue
		}

		if err := field.SetFrom(value); err != nil {
			return t, fmt.Errorf("%w for field %q: %v", ErrSetField, field.Name, err)
		}
	}

	return t, nil
}
