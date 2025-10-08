package environ

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadVariable(t *testing.T) {
	t.Run("returns error when name is empty", func(t *testing.T) {
		variable := Variable[string]{
			Name: "",
			Type: TypeString,
		}
		result, err := loadVariable(variable)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingName)
		assert.Equal(t, "", result)
	})

	t.Run("returns default value when env var does not exist", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		defaultValue := "default"
		variable := Variable[string]{
			Name:    "TEST_VAR",
			Type:    TypeString,
			Default: &defaultValue,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, "default", result)
	})

	t.Run("returns default value when env var is empty string", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		defer os.Unsetenv("TEST_VAR")
		defaultValue := "default"
		variable := Variable[string]{
			Name:    "TEST_VAR",
			Type:    TypeString,
			Default: &defaultValue,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, "default", result)
	})

	t.Run("returns zero value when optional and env var does not exist", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name:     "TEST_VAR",
			Type:     TypeString,
			Optional: true,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("returns zero value when optional and env var is empty", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[int]{
			Name:     "TEST_VAR",
			Type:     TypeInt,
			Optional: true,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("returns error when env var does not exist and not optional", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeString,
		}
		result, err := loadVariable(variable)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingValue)
		assert.Equal(t, "", result)
	})

	t.Run("loads and validates env var successfully", func(t *testing.T) {
		os.Setenv("TEST_VAR", "42")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[int]{
			Name: "TEST_VAR",
			Type: TypeInt,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("returns error for invalid value", func(t *testing.T) {
		os.Setenv("TEST_VAR", "not a number")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[int]{
			Name: "TEST_VAR",
			Type: TypeInt,
		}
		result, err := loadVariable(variable)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInt)
		assert.Equal(t, 0, result)
	})

	t.Run("applies custom validator successfully", func(t *testing.T) {
		os.Setenv("TEST_VAR", "10")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[int]{
			Name: "TEST_VAR",
			Type: TypeInt,
			Validator: func(v int) (int, error) {
				return v * 2, nil
			},
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, 20, result)
	})

	t.Run("returns error when custom validator fails", func(t *testing.T) {
		os.Setenv("TEST_VAR", "10")
		defer os.Unsetenv("TEST_VAR")
		customErr := errors.New("custom validation error")
		variable := Variable[int]{
			Name: "TEST_VAR",
			Type: TypeInt,
			Validator: func(v int) (int, error) {
				return 0, customErr
			},
		}
		result, err := loadVariable(variable)
		assert.Error(t, err)
		assert.ErrorIs(t, err, customErr)
		assert.Equal(t, 0, result)
	})

	t.Run("validates with oneof constraint", func(t *testing.T) {
		os.Setenv("TEST_VAR", "8080")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[int]{
			Name:  "TEST_VAR",
			Type:  TypeInt,
			Oneof: []int{80, 443, 8080},
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, 8080, result)
	})

	t.Run("default value takes precedence over optional", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		defaultValue := "from_default"
		variable := Variable[string]{
			Name:     "TEST_VAR",
			Type:     TypeString,
			Default:  &defaultValue,
			Optional: true,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, "from_default", result)
	})

	t.Run("loads string variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "hello world")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeString,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", result)
	})

	t.Run("loads boolean variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "true")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[bool]{
			Name: "TEST_VAR",
			Type: TypeBoolean,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("loads float variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "3.14")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[float64]{
			Name: "TEST_VAR",
			Type: TypeFloat,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, 3.14, result)
	})

	t.Run("loads port variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "8080")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[int]{
			Name: "TEST_VAR",
			Type: TypePort,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, 8080, result)
	})

	t.Run("loads url variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "https://example.com")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeUrl,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", result)
	})

	t.Run("loads email variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "user@example.com")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeEmail,
		}
		result, err := loadVariable(variable)
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", result)
	})
}

func TestVariableLoad(t *testing.T) {
	t.Run("successfully loads variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "hello")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeString,
		}
		result, err := variable.Load()
		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("wraps error with variable name", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeString,
		}
		result, err := variable.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TEST_VAR")
		assert.ErrorIs(t, err, ErrMissingValue)
		assert.Equal(t, "", result)
	})

	t.Run("loads int variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "42")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[int]{
			Name: "TEST_VAR",
			Type: TypeInt,
		}
		result, err := variable.Load()
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("loads boolean variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "true")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[bool]{
			Name: "TEST_VAR",
			Type: TypeBoolean,
		}
		result, err := variable.Load()
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("loads float variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "3.14")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[float64]{
			Name: "TEST_VAR",
			Type: TypeFloat,
		}
		result, err := variable.Load()
		assert.NoError(t, err)
		assert.Equal(t, 3.14, result)
	})
}

func TestVariableMustLoad(t *testing.T) {
	t.Run("returns value on success", func(t *testing.T) {
		os.Setenv("TEST_VAR", "success")
		defer os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeString,
		}
		result := variable.MustLoad()
		assert.Equal(t, "success", result)
	})

	t.Run("panics on error", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeString,
		}
		assert.Panics(t, func() {
			variable.MustLoad()
		})
	})

	t.Run("panics with correct error", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		variable := Variable[string]{
			Name: "TEST_VAR",
			Type: TypeString,
		}
		assert.PanicsWithError(t, "Err: variable \"TEST_VAR\". Reason: missing required variable", func() {
			variable.MustLoad()
		})
	})
}
