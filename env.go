package confo

import (
	"flag"
	"log"
	"os"
	"strings"
)

var (
	envPrefix = flag.String("env.prefix", "", "Prefix for environment variables.")
)

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
		if v, ok := LookupEnv(fname); ok {
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
	s = strings.ReplaceAll(s, ".", "_")
	return *envPrefix + s
}
