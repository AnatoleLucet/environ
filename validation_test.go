package environ

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInt(t *testing.T) {
	t.Run("parses valid integer", func(t *testing.T) {
		result, err := validateInt("42")
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("parses negative integer", func(t *testing.T) {
		result, err := validateInt("-100")
		assert.NoError(t, err)
		assert.Equal(t, -100, result)
	})

	t.Run("parses zero", func(t *testing.T) {
		result, err := validateInt("0")
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("returns error for invalid integer", func(t *testing.T) {
		result, err := validateInt("not a number")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInt)
		assert.Equal(t, 0, result)
	})

	t.Run("returns error for float string", func(t *testing.T) {
		result, err := validateInt("3.14")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInt)
		assert.Equal(t, 0, result)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		result, err := validateInt("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInt)
		assert.Equal(t, 0, result)
	})
}

func TestValidateFloat(t *testing.T) {
	t.Run("parses valid float", func(t *testing.T) {
		result, err := validateFloat("3.14")
		assert.NoError(t, err)
		assert.Equal(t, 3.14, result)
	})

	t.Run("parses integer as float", func(t *testing.T) {
		result, err := validateFloat("42")
		assert.NoError(t, err)
		assert.Equal(t, 42.0, result)
	})

	t.Run("parses negative float", func(t *testing.T) {
		result, err := validateFloat("-2.5")
		assert.NoError(t, err)
		assert.Equal(t, -2.5, result)
	})

	t.Run("parses zero", func(t *testing.T) {
		result, err := validateFloat("0")
		assert.NoError(t, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("returns error for invalid float", func(t *testing.T) {
		result, err := validateFloat("not a number")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFloat)
		assert.Equal(t, 0.0, result)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		result, err := validateFloat("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFloat)
		assert.Equal(t, 0.0, result)
	})
}

func TestValidateBoolean(t *testing.T) {
	t.Run("parses 'true'", func(t *testing.T) {
		result, err := validateBoolean("true")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("parses 'false'", func(t *testing.T) {
		result, err := validateBoolean("false")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("parses '1'", func(t *testing.T) {
		result, err := validateBoolean("1")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("parses '0'", func(t *testing.T) {
		result, err := validateBoolean("0")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("parses 'on'", func(t *testing.T) {
		result, err := validateBoolean("on")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("parses 'off'", func(t *testing.T) {
		result, err := validateBoolean("off")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("returns error for invalid boolean", func(t *testing.T) {
		result, err := validateBoolean("not a boolean")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidBool)
		assert.False(t, result)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		result, err := validateBoolean("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidBool)
		assert.False(t, result)
	})
}

func TestValidatePort(t *testing.T) {
	t.Run("accepts valid port 80", func(t *testing.T) {
		result, err := validatePort("80")
		assert.NoError(t, err)
		assert.Equal(t, 80, result)
	})

	t.Run("accepts valid port 8080", func(t *testing.T) {
		result, err := validatePort("8080")
		assert.NoError(t, err)
		assert.Equal(t, 8080, result)
	})

	t.Run("accepts minimum port 1", func(t *testing.T) {
		result, err := validatePort("1")
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("accepts maximum port 65535", func(t *testing.T) {
		result, err := validatePort("65535")
		assert.NoError(t, err)
		assert.Equal(t, 65535, result)
	})

	t.Run("rejects port 0", func(t *testing.T) {
		result, err := validatePort("0")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPort)
		assert.Equal(t, 0, result)
	})

	t.Run("rejects negative port", func(t *testing.T) {
		result, err := validatePort("-1")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPort)
		assert.Equal(t, 0, result)
	})

	t.Run("rejects port above 65535", func(t *testing.T) {
		result, err := validatePort("65536")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPort)
		assert.Equal(t, 0, result)
	})

	t.Run("returns error for invalid integer", func(t *testing.T) {
		result, err := validatePort("not a port")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPort)
		assert.Equal(t, 0, result)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		result, err := validatePort("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPort)
		assert.Equal(t, 0, result)
	})
}

func TestValidateUrl(t *testing.T) {
	t.Run("accepts valid http URL", func(t *testing.T) {
		result, err := validateUrl("http://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com", result)
	})

	t.Run("accepts valid https URL", func(t *testing.T) {
		result, err := validateUrl("https://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", result)
	})

	t.Run("accepts URL with path", func(t *testing.T) {
		result, err := validateUrl("https://example.com/path/to/resource")
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/path/to/resource", result)
	})

	t.Run("accepts URL with port", func(t *testing.T) {
		result, err := validateUrl("http://localhost:8080")
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost:8080", result)
	})

	t.Run("accepts URL with query parameters", func(t *testing.T) {
		result, err := validateUrl("https://example.com?foo=bar&baz=qux")
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com?foo=bar&baz=qux", result)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		result, err := validateUrl("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidUrl)
		assert.Equal(t, "", result)
	})

	t.Run("returns error for invalid URL", func(t *testing.T) {
		result, err := validateUrl("not a url")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidUrl)
		assert.Equal(t, "", result)
	})

	t.Run("returns error for URL without protocol", func(t *testing.T) {
		result, err := validateUrl("example.com")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidUrl)
		assert.Equal(t, "", result)
	})
}

func TestValidateEmail(t *testing.T) {
	t.Run("accepts valid email", func(t *testing.T) {
		result, err := validateEmail("user@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", result)
	})

	t.Run("accepts email with subdomain", func(t *testing.T) {
		result, err := validateEmail("user@mail.example.com")
		assert.NoError(t, err)
		assert.Equal(t, "user@mail.example.com", result)
	})

	t.Run("accepts email with plus sign", func(t *testing.T) {
		result, err := validateEmail("user+tag@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "user+tag@example.com", result)
	})

	t.Run("accepts email with dots", func(t *testing.T) {
		result, err := validateEmail("first.last@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "first.last@example.com", result)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		result, err := validateEmail("")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Equal(t, "", result)
	})

	t.Run("returns error for invalid email", func(t *testing.T) {
		result, err := validateEmail("not an email")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Equal(t, "", result)
	})

	t.Run("returns error for email without @", func(t *testing.T) {
		result, err := validateEmail("userexample.com")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Equal(t, "", result)
	})

	t.Run("returns error for email without domain", func(t *testing.T) {
		result, err := validateEmail("user@")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Equal(t, "", result)
	})
}

func TestValidateType(t *testing.T) {
	t.Run("validates string type", func(t *testing.T) {
		result, err := validateType[string](TypeString, "hello")
		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("validates 'str' alias", func(t *testing.T) {
		result, err := validateType[string]("str", "hello")
		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("validates empty type as string", func(t *testing.T) {
		result, err := validateType[string]("", "hello")
		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("validates int type", func(t *testing.T) {
		result, err := validateType[int](TypeInt, "42")
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("validates 'integer' alias", func(t *testing.T) {
		result, err := validateType[int]("integer", "42")
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("validates float type", func(t *testing.T) {
		result, err := validateType[float64](TypeFloat, "3.14")
		assert.NoError(t, err)
		assert.Equal(t, 3.14, result)
	})

	t.Run("validates boolean type", func(t *testing.T) {
		result, err := validateType[bool](TypeBoolean, "true")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("validates 'bool' alias", func(t *testing.T) {
		result, err := validateType[bool]("bool", "false")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("validates port type", func(t *testing.T) {
		result, err := validateType[int](TypePort, "8080")
		assert.NoError(t, err)
		assert.Equal(t, 8080, result)
	})

	t.Run("validates url type", func(t *testing.T) {
		result, err := validateType[string](TypeUrl, "https://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", result)
	})

	t.Run("validates email type", func(t *testing.T) {
		result, err := validateType[string](TypeEmail, "user@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", result)
	})

	t.Run("returns error for unknown type", func(t *testing.T) {
		result, err := validateType[string]("unknown", "value")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUnknownType)
		assert.Equal(t, "", result)
	})
}

func TestValidate(t *testing.T) {
	t.Run("validates value with correct type", func(t *testing.T) {
		variable := Variable[int]{
			Name: "TEST",
			Type: TypeInt,
		}
		result, err := validate(variable, "42")
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("validates value with oneof constraint", func(t *testing.T) {
		variable := Variable[int]{
			Name:  "PORT",
			Type:  TypeInt,
			Oneof: []int{80, 443, 8080},
		}
		result, err := validate(variable, "8080")
		assert.NoError(t, err)
		assert.Equal(t, 8080, result)
	})

	t.Run("returns error when value not in oneof", func(t *testing.T) {
		variable := Variable[int]{
			Name:  "PORT",
			Type:  TypeInt,
			Oneof: []int{80, 443, 8080},
		}
		result, err := validate(variable, "3000")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotInOneof)
		assert.Equal(t, 0, result)
	})

	t.Run("validates string with oneof constraint", func(t *testing.T) {
		variable := Variable[string]{
			Name:  "ENV",
			Type:  TypeString,
			Oneof: []string{"dev", "staging", "prod"},
		}
		result, err := validate(variable, "prod")
		assert.NoError(t, err)
		assert.Equal(t, "prod", result)
	})

	t.Run("returns error for invalid type", func(t *testing.T) {
		variable := Variable[int]{
			Name: "TEST",
			Type: TypeInt,
		}
		result, err := validate(variable, "not a number")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInt)
		assert.Equal(t, 0, result)
	})
}
