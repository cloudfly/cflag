package confo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func getFileWithEnv(file, env string) (string, time.Time, error) {
	if env == "" {
		return "", time.Now(), fmt.Errorf("no env defined")
	}
	var (
		envFile string
		extname = path.Ext(file)
	)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, fileInfo.ModTime(), nil
	}
	return "", time.Now(), fmt.Errorf("no env file %s found", file)
}

func parseArgument() map[string]string {
	args := os.Args[1:]
	result := map[string]string{}
	idx := 1
	currentName := ""
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			// find new argument name
			if currentName != "" {
				result[currentName] = "true"
			}
			currentName = strings.TrimLeft(args[i], "-")
			//  --your-argname=argvalue
			if name, value, ok := strings.Cut(currentName, "="); ok {
				result[name] = value
				currentName = ""
			}
		} else {
			if currentName != "" {
				result[currentName] = args[i]
				currentName = ""
			} else {
				result[fmt.Sprintf("$%d", idx)] = args[i]
				idx++
			}
		}
	}
	if currentName != "" {
		result[currentName] = "true"
	}
	return result
}

func processFile(config interface{}, file string, errorOnUnmatchedKeys bool) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}

	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		return yaml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".toml"):
		return unmarshalToml(data, config)
	case strings.HasSuffix(file, ".json"):
		return unmarshalJSON(data, config)
	default:
		if data[0] == '{' { // json
			return unmarshalJSON(data, config)
		}
		yamlError := yaml.Unmarshal(data, config)
		if yamlError == nil {
			return nil
		}
		tomlError := unmarshalToml(data, config)
		if tomlError == nil {
			return nil
		}
		return fmt.Errorf("decode config error: %w", errors.Join(yamlError, tomlError))
	}
}

func unmarshalToml(data []byte, config interface{}) error {
	_, err := toml.Decode(string(data), config)
	return err
}

// unmarshalJSON unmarshals the given data into the config interface.
func unmarshalJSON(data []byte, config interface{}) error {
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(config)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func getEnvPrefixForStruct(prefixes []string, fieldStruct *reflect.StructField) []string {
	envName := fieldStruct.Tag.Get("env")
	if envName == "-" {
		return prefixes
	} else if envName == "" {
		return append(prefixes, fieldStruct.Name)
	}
	return append(prefixes, envName)
}

func processDefault(config interface{}) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		switch field.Kind() {
		case reflect.Struct:
			if err := processDefault(field.Addr().Interface()); err != nil {
				return err
			}
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := processDefault(field.Index(i).Addr().Interface()); err != nil {
						return err
					}
				}
			}
		default:
			if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank {
				// Set default configuration if blank
				if value := fieldStruct.Tag.Get("default"); value != "" {
					if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func processRequired(config any) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		fieldStruct := configType.Field(i)
		field := configValue.Field(i)

		if !field.CanInterface() {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			if err := processRequired(field.Addr().Interface()); err != nil {
				return err
			}
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := processRequired(field.Index(i).Addr().Interface()); err != nil {
						return err
					}
				}
			}
		default:
			if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank && fieldStruct.Tag.Get("required") == "true" {
				// return error if it is required but blank
				return errors.New(fieldStruct.Name + " is required, but empty")
			}
		}

	}
	return nil
}

func processEnv(config interface{}, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			envNames    []string
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
			envName     = fieldStruct.Tag.Get("env") // read configuration from shell env
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		if envName == "" {
			envNames = append(envNames, strings.Join(append(prefixes, fieldStruct.Name), "_"))                  // Confo_Db_Name
			envNames = append(envNames, strings.ToUpper(strings.Join(append(prefixes, fieldStruct.Name), "_"))) // CONFO_DB_NAME
		} else {
			for _, name := range strings.Split(envName, ",") {
				name := name
				envNames = append(envNames, name)
			}
		}

		// confo.c.LogFn(nil, "Trying to load struct `%v`'s field `%v` from env %v\n", configType.Name(), fieldStruct.Name, strings.Join(envNames, ", "))

		// Load From Shell ENV
		for _, env := range envNames {
			if value := os.Getenv(env); value != "" {
				// confo.c.LogFn(nil, "Loading configuration for struct `%v`'s field `%v` from env %v...\n", configType.Name(), fieldStruct.Name, env)
				switch reflect.Indirect(field).Kind() {
				case reflect.Bool:
					switch strings.ToLower(value) {
					case "", "0", "false":
						field.Set(reflect.ValueOf(false))
					default:
						field.Set(reflect.ValueOf(true))
					}
				case reflect.String:
					field.Set(reflect.ValueOf(value))
				default:
					if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
						return err
					}
				}
				break
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := processEnv(field.Addr().Interface(), getEnvPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			if arrLen := field.Len(); arrLen > 0 {
				for i := 0; i < arrLen; i++ {
					if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
						if err := processEnv(field.Index(i).Addr().Interface(), append(getEnvPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(i))...); err != nil {
							return err
						}
					}
				}
			} else {
				// load slice from env
				newVal := reflect.New(field.Type().Elem()).Elem()
				if newVal.Kind() == reflect.Struct {
					idx := 0
					for {
						newVal = reflect.New(field.Type().Elem()).Elem()
						if err := processEnv(newVal.Addr().Interface(), append(getEnvPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(idx))...); err != nil {
							return err
						} else if reflect.DeepEqual(newVal.Interface(), reflect.New(field.Type().Elem()).Elem().Interface()) {
							break
						} else {
							idx++
							field.Set(reflect.Append(field, newVal))
						}
					}
				}
			}
		}
	}
	return nil
}

func processArgs(config any, prefix ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	configType := configValue.Type()
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}
	argMap := parseArgument()

	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		switch field.Kind() {
		case reflect.Struct:
			if err := processArgs(field.Addr().Interface()); err != nil {
				return err
			}
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := processArgs(field.Index(i).Addr().Interface()); err != nil {
						return err
					}
				}
			}
		default:
			var argNames []string
			if s := fieldStruct.Tag.Get("arg"); s != "" {
				for _, name := range strings.Split(s, ",") {
					name := name
					argNames = append(argNames, strings.Join(append(prefix, name), "-"))
				}
			} else {
				argNames = append(argNames, strings.Join(append(prefix, fieldStruct.Name), "-"))
			}
			for _, name := range argNames {
				value := argMap[name]
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
