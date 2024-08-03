package confo

import (
	"os"
	"strings"
)

// LookupEnv returns the expanded environment variable value for the given name.
//
// The expanded means that `%{ENV_VAR}` placeholders in env var value are replaced
// with the corresponding ENV_VAR values (recursively).
//
// false is returned if environment variable isn't found.
func LookupEnv(name string) (string, bool) {
	value, ok := envVars[name]
	return value, ok
}

var envVars = func() map[string]string {
	envs := os.Environ()
	return parseEnvVars(envs)
}()

func parseEnvVars(envs []string) map[string]string {
	m := make(map[string]string, len(envs))
	for _, env := range envs {
		n := strings.IndexByte(env, '=')
		if n < 0 {
			m[env] = ""
			continue
		}
		name := env[:n]
		value := env[n+1:]
		m[name] = value
	}
	return m
}
