# Confo

Golang Configuration tool that support YAML, JSON, TOML, Shell Environment, Command Line (Supports Go 1.19+)

## Usage

```go
package main

import (
	"fmt"
	"github.com/cloudfly/confo"
)

var Config = struct {
	APPName string `default:"myapp"`
	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}
	Token string `env:"TOKEN" arg:"-"`
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

```

## No file, loading from environment and command line

```go
func main() {
	confo.Load(&Config)
	fmt.Printf("config: %#v", Config)
}
```

```bash
# config the token by environment for Config.Token
export TOKEN='your_token'
```

## Priority Order

**Command line > Environment > Last file in Load() argument ... > First file in Load() argument**

eg. 
```go

var Config = struct {
	APPName string `default:"myapp"`
	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}
	Token string `env:"TOKEN" arg:"-"`
}{}

c := confo.New(&confo.Config{Env: "prod"}) // setting the env
c.Load(&Config, "common.yml", "config.yml")
```

Below is Confo's loading order, **back loading operation will overwrite previous.**
1. loading default value defined in structure's tag
1. loading common.yml
1. loading common.prod.yml if exist
2. loading config.yml
3. loading config.prod.yml if exist
4. loading environment
5. loading arguments in command line 

## Auto Reload Mode

Confo can auto reload configuration based on time

```go
// auto reload configuration every second
confo.New(&confo.Config{AutoReload: true}).Load(&Config, "config.json")

// auto reload configuration every minute
confo.New(&confo.Config{AutoReload: true, AutoReloadInterval: time.Minute}).Load(&Config, "config.json")
```

Auto Reload Callback

```go
confo.New(&confo.Config{AutoReload: true, AutoReloadCallback: func(config any) {
    fmt.Printf("%v changed", config)
}}).Load(&Config, "config.json")
```

## Contributing

You can help to make the project better.

## Author

**cloudfly**

* <http://github.com/cloudfly>
* <chenyunfei.cs@gmail.com>

## License

Released under the MIT License
