package environ

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("loads struct with single string field", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string"`
		}
		os.Setenv("APP_NAME", "my-app")
		defer os.Unsetenv("APP_NAME")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "my-app", result.Name)
	})

	t.Run("loads struct with multiple fields", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string"`
			Port int    `env:"name=APP_PORT, type=int"`
		}
		os.Setenv("APP_NAME", "my-app")
		os.Setenv("APP_PORT", "8080")
		defer os.Unsetenv("APP_NAME")
		defer os.Unsetenv("APP_PORT")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "my-app", result.Name)
		assert.Equal(t, 8080, result.Port)
	})

	t.Run("loads struct with all types", func(t *testing.T) {
		type Config struct {
			Str   string  `env:"name=STR, type=string"`
			Num   int     `env:"name=NUM, type=int"`
			Float float64 `env:"name=FLOAT, type=float"`
			Bool  bool    `env:"name=BOOL, type=bool"`
			Port  int     `env:"name=PORT, type=port"`
			Url   string  `env:"name=URL, type=url"`
			Email string  `env:"name=EMAIL, type=email"`
		}
		os.Setenv("STR", "hello")
		os.Setenv("NUM", "42")
		os.Setenv("FLOAT", "3.14")
		os.Setenv("BOOL", "true")
		os.Setenv("PORT", "8080")
		os.Setenv("URL", "https://example.com")
		os.Setenv("EMAIL", "user@example.com")
		defer func() {
			os.Unsetenv("STR")
			os.Unsetenv("NUM")
			os.Unsetenv("FLOAT")
			os.Unsetenv("BOOL")
			os.Unsetenv("PORT")
			os.Unsetenv("URL")
			os.Unsetenv("EMAIL")
		}()

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "hello", result.Str)
		assert.Equal(t, 42, result.Num)
		assert.Equal(t, 3.14, result.Float)
		assert.True(t, result.Bool)
		assert.Equal(t, 8080, result.Port)
		assert.Equal(t, "https://example.com", result.Url)
		assert.Equal(t, "user@example.com", result.Email)
	})

	t.Run("uses default values when env vars not set", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string, default=default-app"`
			Port int    `env:"name=APP_PORT, type=int, default=3000"`
		}
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_PORT")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "default-app", result.Name)
		assert.Equal(t, 3000, result.Port)
	})

	t.Run("loads optional fields with zero values", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string, optional"`
			Port int    `env:"name=APP_PORT, type=int, optional"`
		}
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_PORT")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "", result.Name)
		assert.Equal(t, 0, result.Port)
	})

	t.Run("returns error for missing required field", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string"`
		}
		os.Unsetenv("APP_NAME")

		result, err := Load[Config]()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingValue)
		assert.Equal(t, "", result.Name)
	})

	t.Run("returns error for invalid value", func(t *testing.T) {
		type Config struct {
			Port int `env:"name=APP_PORT, type=int"`
		}
		os.Setenv("APP_PORT", "not-a-number")
		defer os.Unsetenv("APP_PORT")

		result, err := Load[Config]()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidInt)
		assert.Equal(t, 0, result.Port)
	})

	t.Run("validates oneof constraint", func(t *testing.T) {
		type Config struct {
			Env string `env:"name=APP_ENV, type=string, oneof=dev|staging|prod"`
		}
		os.Setenv("APP_ENV", "prod")
		defer os.Unsetenv("APP_ENV")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "prod", result.Env)
	})

	t.Run("returns error for invalid oneof value", func(t *testing.T) {
		type Config struct {
			Env string `env:"name=APP_ENV, type=string, oneof=dev|staging|prod"`
		}
		os.Setenv("APP_ENV", "test")
		defer os.Unsetenv("APP_ENV")

		result, err := Load[Config]()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotInOneof)
		assert.Equal(t, "", result.Env)
	})

	t.Run("loads empty struct", func(t *testing.T) {
		type Config struct{}

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, Config{}, result)
	})

	t.Run("ignores fields without env tag", func(t *testing.T) {
		type Config struct {
			Name   string `env:"name=APP_NAME, type=string"`
			Ignore string
		}
		os.Setenv("APP_NAME", "my-app")
		defer os.Unsetenv("APP_NAME")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "my-app", result.Name)
		assert.Equal(t, "", result.Ignore)
	})

	t.Run("loads with description tag", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string, desc=Application name"`
		}
		os.Setenv("APP_NAME", "my-app")
		defer os.Unsetenv("APP_NAME")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "my-app", result.Name)
	})

	t.Run("env var overrides default", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string, default=default-name"`
		}
		os.Setenv("APP_NAME", "custom-name")
		defer os.Unsetenv("APP_NAME")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "custom-name", result.Name)
	})

	t.Run("loads boolean with different formats", func(t *testing.T) {
		type Config struct {
			Flag1 bool `env:"name=FLAG1, type=bool"`
			Flag2 bool `env:"name=FLAG2, type=bool"`
			Flag3 bool `env:"name=FLAG3, type=bool"`
		}
		os.Setenv("FLAG1", "true")
		os.Setenv("FLAG2", "1")
		os.Setenv("FLAG3", "on")
		defer func() {
			os.Unsetenv("FLAG1")
			os.Unsetenv("FLAG2")
			os.Unsetenv("FLAG3")
		}()

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.True(t, result.Flag1)
		assert.True(t, result.Flag2)
		assert.True(t, result.Flag3)
	})

	t.Run("validates port range", func(t *testing.T) {
		type Config struct {
			Port int `env:"name=APP_PORT, type=port"`
		}
		os.Setenv("APP_PORT", "99999")
		defer os.Unsetenv("APP_PORT")

		result, err := Load[Config]()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPort)
		assert.Equal(t, 0, result.Port)
	})

	t.Run("validates url format", func(t *testing.T) {
		type Config struct {
			Url string `env:"name=APP_URL, type=url"`
		}
		os.Setenv("APP_URL", "not-a-url")
		defer os.Unsetenv("APP_URL")

		result, err := Load[Config]()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidUrl)
		assert.Equal(t, "", result.Url)
	})

	t.Run("validates email format", func(t *testing.T) {
		type Config struct {
			Email string `env:"name=APP_EMAIL, type=email"`
		}
		os.Setenv("APP_EMAIL", "not-an-email")
		defer os.Unsetenv("APP_EMAIL")

		result, err := Load[Config]()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Equal(t, "", result.Email)
	})

	t.Run("uses type aliases", func(t *testing.T) {
		type Config struct {
			Name  string `env:"name=NAME, type=str"`
			Count int    `env:"name=COUNT, type=integer"`
			Flag  bool   `env:"name=FLAG, type=bool"`
		}
		os.Setenv("NAME", "test")
		os.Setenv("COUNT", "10")
		os.Setenv("FLAG", "true")
		defer func() {
			os.Unsetenv("NAME")
			os.Unsetenv("COUNT")
			os.Unsetenv("FLAG")
		}()

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 10, result.Count)
		assert.True(t, result.Flag)
	})
}

func TestMustLoad(t *testing.T) {
	t.Run("returns value on success", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string"`
		}
		os.Setenv("APP_NAME", "my-app")
		defer os.Unsetenv("APP_NAME")

		result := MustLoad[Config]()
		assert.Equal(t, "my-app", result.Name)
	})

	t.Run("panics on error", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string"`
		}
		os.Unsetenv("APP_NAME")

		assert.Panics(t, func() {
			MustLoad[Config]()
		})
	})

	t.Run("panics with error containing variable name", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string"`
		}
		os.Unsetenv("APP_NAME")

		defer func() {
			r := recover()
			assert.NotNil(t, r)
			err, ok := r.(error)
			assert.True(t, ok)
			assert.Contains(t, err.Error(), "APP_NAME")
		}()

		MustLoad[Config]()
	})
}

func TestLoadEdgeCases(t *testing.T) {
	t.Run("handles empty env var with default", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string, default=fallback"`
		}
		os.Setenv("APP_NAME", "")
		defer os.Unsetenv("APP_NAME")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "fallback", result.Name)
	})

	t.Run("handles whitespace in env var", func(t *testing.T) {
		type Config struct {
			Name string `env:"name=APP_NAME, type=string"`
		}
		os.Setenv("APP_NAME", "  spaces  ")
		defer os.Unsetenv("APP_NAME")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "  spaces  ", result.Name)
	})

	t.Run("handles negative numbers", func(t *testing.T) {
		type Config struct {
			Value int `env:"name=VALUE, type=int"`
		}
		os.Setenv("VALUE", "-42")
		defer os.Unsetenv("VALUE")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, -42, result.Value)
	})

	t.Run("handles zero as valid value", func(t *testing.T) {
		type Config struct {
			Count int `env:"name=COUNT, type=int"`
		}
		os.Setenv("COUNT", "0")
		defer os.Unsetenv("COUNT")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, 0, result.Count)
	})

	t.Run("handles false as valid value", func(t *testing.T) {
		type Config struct {
			Flag bool `env:"name=FLAG, type=bool"`
		}
		os.Setenv("FLAG", "false")
		defer os.Unsetenv("FLAG")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.False(t, result.Flag)
	})

	t.Run("loads struct with mix of set and unset optional fields", func(t *testing.T) {
		type Config struct {
			Name  string `env:"name=NAME, type=string, optional"`
			Email string `env:"name=EMAIL, type=string, optional"`
		}
		os.Setenv("NAME", "test")
		os.Unsetenv("EMAIL")
		defer os.Unsetenv("NAME")

		result, err := Load[Config]()
		assert.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, "", result.Email)
	})
}
