package confo

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

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

var NewBool = flag.Bool
var NewString = flag.String
var NewInt = flag.Int
var NewInt64 = flag.Int64
var NewFloat = flag.Float64
