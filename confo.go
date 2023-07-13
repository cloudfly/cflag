package confo

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

var (
	// Default confo instance
	Default *Confo
)

func init() {
	var c Config
	if err := New(nil).Load(&c); err != nil {
		panic(fmt.Errorf("init default Confo error: %w", err))
	}
	Default = New(&c)
}

type Confo struct {
	c        *Config
	modTimes map[string]time.Time
}

type Config struct {
	// Environment name, Confo will auto try to load the file <filename>.<env>.<ext>。
	// eg. you Load a file config.yml and the this Env == testing，then confo will try to load file config.testing.yml in the directory of config.yml
	Env string `env:"CONFO_ENV" arg:"confo-env"`
	// ArgPrefix set the common prefix for argument name
	ArgPrefix string `env:"CONFO_ARG_PREFIX" arg:"confo-arg-prefix"`
	// EnvPrefix set the common prefix for environment name
	EnvPrefix string `env:"CONFO_ENV_PREFIX" arg:"confo-env-prefix"`
	// set the log func for confo, default is nil, which means drop all the message
	LogFn func(err error, format string, v ...any)
	// enable auto reload config files, default false
	AutoReload bool `env:"CONFO_AUTO_RELOAD" arg:"confo-auto-reload"`
	// time interval for checking config files, default 3 seconds.
	AutoReloadInterval time.Duration `env:"CONFO_AUTO_RELOAD_INTERVAL" arg:"confo-auto-reload-interval"`
	// the callback function after reloading the config
	AutoReloadCallback func(config any)
}

// New initialize a Confo instance
func New(c *Config) *Confo {
	if c == nil {
		c = &Config{}
	}

	if c.AutoReloadInterval == 0 {
		c.AutoReloadInterval = time.Second * 3
	}

	if c.LogFn == nil {
		c.LogFn = func(err error, format string, v ...any) {}
	}

	return &Confo{c: c, modTimes: make(map[string]time.Time)}
}

// Load will unmarshal configurations to struct from files that you provide one by one
//
// configuration in later file will overwrite the previous, so the last file has top priority
func (confo *Confo) Load(target interface{}, files ...string) (err error) {
	defaultValue := reflect.Indirect(reflect.ValueOf(target))
	if !defaultValue.CanAddr() {
		return fmt.Errorf("Config %v should be addressable", target)
	}
	err, _ = confo.load(target, false, files...)

	if confo.c.AutoReload {
		go confo.autoLoad(target, files...)
	}
	return
}

func (confo *Confo) autoLoad(target any, files ...string) {
	defaultValue := reflect.Indirect(reflect.ValueOf(target))
	ticker := time.NewTicker(confo.c.AutoReloadInterval)
	defer ticker.Stop()
	for range ticker.C {
		ptr := reflect.New(reflect.ValueOf(target).Elem().Type())
		ptr.Elem().Set(defaultValue)

		if err, changed := confo.load(ptr.Interface(), true, files...); err == nil && changed {
			confo.c.LogFn(err, "Config files changed, reload it")
			reflect.ValueOf(target).Elem().Set(ptr.Elem())
			if confo.c.AutoReloadCallback != nil {
				confo.c.AutoReloadCallback(target)
			}
		} else if err != nil {
			confo.c.LogFn(err, "Failed to reload configuration from %v", files)
		}
	}
}
func (confo *Confo) getFiles(files ...string) ([]string, map[string]time.Time) {
	var resultFiles []string
	var modTimeMap = map[string]time.Time{}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]

		// check configuration
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			resultFiles = append(resultFiles, file)
			modTimeMap[file] = fileInfo.ModTime()
		}

		// check configuration with env
		if file, modTime, err := getFileWithEnv(file, confo.c.Env); err == nil {
			resultFiles = append(resultFiles, file)
			modTimeMap[file] = modTime
		}
	}
	return resultFiles, modTimeMap
}

func (confo *Confo) load(config interface{}, watchMode bool, files ...string) (err error, changed bool) {
	defer confo.c.LogFn(err, "Configuration: %#v", config)

	configFiles, configModTimeMap := confo.getFiles(files...)

	if watchMode {
		if len(configModTimeMap) == len(confo.modTimes) {
			var changed bool
			for f, t := range configModTimeMap {
				if v, ok := confo.modTimes[f]; !ok || t.After(v) {
					changed = true
				}
			}
			if !changed {
				return nil, false
			}
		}
	}

	for _, file := range configFiles {
		confo.c.LogFn(nil, "Loading configurations from file '%v'...\n", file)
		if err = processFile(config, file, false); err != nil {
			return err, true
		}
	}

	confo.modTimes = configModTimeMap

	if confo.c.EnvPrefix == "" || confo.c.EnvPrefix == "-" {
		err = processEnv(config)
	} else {
		err = processEnv(config, confo.c.EnvPrefix)
	}

	if confo.c.ArgPrefix == "" || confo.c.ArgPrefix == "-" {
		err = processArgs(config)
	} else {
		err = processArgs(config, confo.c.ArgPrefix)
	}

	// process defaults
	processDefault(config)

	if err := processRequired(config); err != nil {
		return err, false
	}

	return err, true
}

// GetConfig return a copy of Config of Confo
func (confo *Confo) GetConfig() *Config {
	c2 := *confo.c
	return &c2
}

// Load will unmarshal configurations to struct from files that you provide by using Default Confo
func Load(config interface{}, files ...string) error {
	return Default.Load(config, files...)
}
