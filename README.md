<h1 align="center"><code>environ</code></h1>

<p align="center">Load environement variables without the hustle.</p>

```go
type Envs struct {
	Url    string `env:"name=URL, type=url, default=http://localhost"`
	Port   int    `env:"name=PORT, type=port, default=8080"`
}

envs, err := environ.Load[Envs]()

envs.Url  // http://localhost
envs.Port // 8080
```

## Installation

```bash
go get github.com/AnatoleLucet/environ
```

## Usage

With `environ`, you can load environments variables in three ways.

> There's no good or bad option here. `environ` wants to provide a way to load environment variables for you as simply and flexibly as possible in any situation, without any compromise. So pick the one that makes most sense for you!

#### Tag-based loading

```go
type Envs struct {
    Url      string `env:"name=URL, type=url, default=http://localhost"`
    Port     int    `env:"name=PORT, type=port, optional"`
    Secret   string `env:"name=SECRET, type=string"`
}

func LoadEnvs() {
    envs, err := environ.Load[Envs]()
    // or with panic:
    // envs := environ.MustLoad[Envs]()
}
```

#### Builder-based loading

```go
func LoadEnvs() {
    url, err := environ.Url("URL").Default("http://localhost").Load()
    port, err := environ.Port("PORT").Optional().Load()
    secret, err := environ.String("SECRET").Load()
    // or with panic:
    // secret := environ.String("SECRET").MustLoad()
}
```

#### Struct-based loading

```go
func LoadEnvs() {
    defaultUrl := "http://localhost"

    url, err := environ.Variable[string]{
        Name: "URL",
        Type: environ.TypeUrl,
        Default: &defaultUrl,
    }.Load()
    port, err := environ.Variable[int]{
        Name: "PORT",
        Type: environ.TypePort,
        Optional: true,
    }.Load()
    secret, err := environ.Variable[string]{
        Name: "SECRET",
        Type: environ.TypeString,
    }.Load()
    // or with panic:
    // secret := environ.Variable[string]{Name: "SECRET", Type: environ.TypeString}.MustLoad()
}
```

### Options

| Name       | Go Type                           | Description                                     | Tag Example                                          |
| ---------- | --------------------------------- | ----------------------------------------------- | ---------------------------------------------------- |
| `name`     | `string`                          | Name of the environment variable to load        | `env="name=URL"`                                     |
| `type`     | [`VariableType`](#Variable-Types) | Expected type of the variable's value           | `env="type=bool"`                                    |
| `default`  | `T`                               | Fallback value if the variable is not defined   | `env="default=http://localhost"`                     |
| `optional` | `bool`                            | If the variable is allowed to be empty of not   | `env="optional"`                                     |
| `desc`     | `string`                          | Description of the variable                     | `env="desc=A short description about this variable"` |
| `oneof`    | `[]T`                             | Allow-list of values the variable can be set to | `env="oneof=80\|3000\|8080"`                         |

### Variable Types

| Type          | Go Type   | Description                                         | Tag Example         |
| ------------- | --------- | --------------------------------------------------- | ------------------- |
| `TypeString`  | `string`  | Any string value                                    | `env:"type=string"` |
| `TypeInt`     | `int`     | Integer number                                      | `env:"type=int"`    |
| `TypeFloat`   | `float64` | Floating point number                               | `env:"type=float"`  |
| `TypeBoolean` | `bool`    | Boolean value (accepts: true, false, 0, 1, on, off) | `env:"type=bool"`   |
| `TypePort`    | `int`     | Valid TCP port number (1-65535)                     | `env:"type=port"`   |
| `TypeUrl`     | `string`  | Valid URL with protocol and hostname                | `env:"type=url"`    |
| `TypeEmail`   | `string`  | Valid email address                                 | `env:"type=email"`  |
