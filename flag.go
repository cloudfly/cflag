package confo

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	envPrefix = flag.String("env.prefix", "", "Prefix for environment variables.")
)

var NewBool = flag.Bool
var NewString = flag.String
var NewInt = flag.Int
var NewInt64 = flag.Int64
var NewFloat = flag.Float64

// WriteFlags writes all the explicitly set flags to w.
func WriteFlags(w io.Writer) {
	flag.Visit(func(f *flag.Flag) {
		lname := strings.ToLower(f.Name)
		value := f.Value.String()
		if IsSecretFlag(lname) {
			value = "secret"
		}
		fmt.Fprintf(w, "-%s=%q\n", f.Name, value)
	})
}

// Visit all the flag name and values
func Visit(fn func(string, string)) {
	flag.Visit(func(f *flag.Flag) {
		lname := strings.ToLower(f.Name)
		value := f.Value.String()
		if IsSecretFlag(lname) {
			value = "secret"
		}
		fn(lname, value)
	})
}

// Parse parses environment vars and command-line flags.
//
// Flags set via command-line override flags set via environment vars.
//
// This function must be called instead of flag.Parse() before using any flags in the program.
func Parse() {
	ParseFlagSet(flag.CommandLine, os.Args[1:])
}

// ParseFlagSet parses the given args into the given fs.
func ParseFlagSet(fs *flag.FlagSet, args []string) {
	if err := fs.Parse(args); err != nil {
		// Do not use lib/logger here, since it is uninitialized yet.
		log.Fatalf("cannot parse flags %q: %s", args, err)
	}

	// Remember explicitly set command-line flags.
	flagsSet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})

	// Obtain the remaining flag values from environment vars.
	fs.VisitAll(func(f *flag.Flag) {
		if flagsSet[f.Name] {
			// The flag is explicitly set via command-line.
			return
		}
		// Get flag value from environment var.
		fname := getEnvFlagName(f.Name)
		if v := os.Getenv(fname); v != "" {
			if err := fs.Set(f.Name, v); err != nil {
				// Do not use lib/logger here, since it is uninitialized yet.
				log.Fatalf("cannot set flag %s to %q, which is read from env var %q: %s", f.Name, v, fname, err)
			}
		}
	})
}

func getEnvFlagName(s string) string {
	// Substitute dots with underscores, since env var names cannot contain dots.
	// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/311#issuecomment-586354129 for details.
	return strings.ToUpper(*envPrefix + strings.ReplaceAll(s, ".", "_"))
}

func envHelp(s string) string {
	return fmt.Sprintf("(env: %s)", getEnvFlagName(s))
}
