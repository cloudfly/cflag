# Confo

Golang Configuration tool that support YAML, JSON, TOML, Shell Environment (Supports Go 1.19+)

## Usage

```go
package main

import (
	"fmt"
	"github.com/cloudfly/confo"
)

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

func main() {
	confo.Load(&Config, "config.yml")
	fmt.Printf("config: %#v", Config)
}
```

With configuration file *config.yml*:

```yaml
appname: test

db:
    name:     test
    user:     test
    password: test
    port:     1234

contacts:
- name: i test
  email: test@test.com
```

## Auto Reload Mode

Confo can auto reload configuration based on time

```go
// auto reload configuration every second
Confo.New(&Confo.Config{AutoReload: true}).Load(&Config, "config.json")

// auto reload configuration every minute
Confo.New(&Confo.Config{AutoReload: true, AutoReloadInterval: time.Minute}).Load(&Config, "config.json")
```

Auto Reload Callback

```go
Confo.New(&Confo.Config{AutoReload: true, AutoReloadCallback: func(config interface{}) {
    fmt.Printf("%v changed", config)
}}).Load(&Config, "config.json")
```

# Advanced Usage

* Load mutiple configurations

```go
// Earlier configurations have higher priority
Confo.Load(&Config, "application.yml", "database.json")
```

* Return error on unmatched keys

Return an error on finding keys in the config file that do not match any fields in the config struct.
In the example below, an error will be returned if config.toml contains keys that do not match any fields in the ConfigStruct struct.
If ErrorOnUnmatchedKeys is not set, it defaults to false.

Note that for json files, setting ErrorOnUnmatchedKeys to true will have an effect only if using go 1.10 or later.

```go
err := Confo.New(&Confo.Config{ErrorOnUnmatchedKeys: true}).Load(&ConfigStruct, "config.toml")
```

* Load configuration by environment

Use `Confo_ENV` to set environment, if `Confo_ENV` not set, environment will be `development` by default, and it will be `test` when running tests with `go test`

```go
// config.go
Confo.Load(&Config, "config.json")

$ go run config.go
// Will load `config.json`, `config.development.json` if it exists
// `config.development.json` will overwrite `config.json`'s configuration
// You could use this to share same configuration across different environments

$ Confo_ENV=production go run config.go
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration

$ go test
// Will load `config.json`, `config.test.json` if it exists
// `config.test.json` will overwrite `config.json`'s configuration

$ Confo_ENV=production go test
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration
```

```go
// Set environment by config
Confo.New(&Confo.Config{Environment: "production"}).Load(&Config, "config.json")
```

* Example Configuration

```go
// config.go
Confo.Load(&Config, "config.yml")

$ go run config.go
// Will load `config.example.yml` automatically if `config.yml` not found and print warning message
```

* Load From Shell Environment

```go
$ Confo_APPNAME="hello world" Confo_DB_NAME="hello world" go run config.go
// Load configuration from shell environment, it's name is {{prefix}}_FieldName
```

```go
// You could overwrite the prefix with environment Confo_ENV_PREFIX, for example:
$ Confo_ENV_PREFIX="WEB" WEB_APPNAME="hello world" WEB_DB_NAME="hello world" go run config.go

// Set prefix by config
confo.New(&confo.Config{EnvPrefix: "WEB"}).Load(&Config, "config.json")
```

* With flags

```go
func main() {
	config := flag.String("file", "config.yml", "configuration file")
	flag.StringVar(&Config.APPName, "name", "", "app name")
	flag.StringVar(&Config.DB.Name, "db-name", "", "database name")
	flag.StringVar(&Config.DB.User, "db-user", "root", "database user")
	flag.Parse()

	os.Setenv("Confo_ENV_PREFIX", "-")
	Confo.Load(&Config, *config)
	// Confo.Load(&Config) // only load configurations from shell env & flag
}
```

## Contributing

You can help to make the project better.

## Author

**cloudfly**

* <http://github.com/cloudfly>
* <chenyunfei.cs@gmail.com>

## License

Released under the MIT License
