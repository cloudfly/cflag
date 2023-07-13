package confo

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Bytes is a value for holding size in bytes.
//
// It supports the following optional suffixes for values: KB, MB, GB, KiB, MiB, GiB.
type Bytes struct {
	// N contains parsed value for the given flag.
	N int

	valueString string
}

// String implements flag.Value interface
func (b *Bytes) String() string {
	return b.valueString
}

// Set implements flag.Value interface
func (b *Bytes) Set(value string) error {
	value = normalizeBytesString(value)
	switch {
	case strings.HasSuffix(value, "KB"):
		f, err := strconv.ParseFloat(value[:len(value)-2], 64)
		if err != nil {
			return err
		}
		b.N = int(f * 1000)
		b.valueString = value
		return nil
	case strings.HasSuffix(value, "MB"):
		f, err := strconv.ParseFloat(value[:len(value)-2], 64)
		if err != nil {
			return err
		}
		b.N = int(f * 1000 * 1000)
		b.valueString = value
		return nil
	case strings.HasSuffix(value, "GB"):
		f, err := strconv.ParseFloat(value[:len(value)-2], 64)
		if err != nil {
			return err
		}
		b.N = int(f * 1000 * 1000 * 1000)
		b.valueString = value
		return nil
	case strings.HasSuffix(value, "KiB"):
		f, err := strconv.ParseFloat(value[:len(value)-3], 64)
		if err != nil {
			return err
		}
		b.N = int(f * 1024)
		b.valueString = value
		return nil
	case strings.HasSuffix(value, "MiB"):
		f, err := strconv.ParseFloat(value[:len(value)-3], 64)
		if err != nil {
			return err
		}
		b.N = int(f * 1024 * 1024)
		b.valueString = value
		return nil
	case strings.HasSuffix(value, "GiB"):
		f, err := strconv.ParseFloat(value[:len(value)-3], 64)
		if err != nil {
			return err
		}
		b.N = int(f * 1024 * 1024 * 1024)
		b.valueString = value
		return nil
	default:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		b.N = int(f)
		b.valueString = value
		return nil
	}
}

func (b *Bytes) UnmarshalJSON(value []byte) error {
	return b.Set(string(value))
}

func (b *Bytes) UnmarshalYAML(node *yaml.Node) error {
	return b.Set(node.Value)
}

func normalizeBytesString(s string) string {
	s = strings.ToUpper(s)
	return strings.ReplaceAll(s, "I", "i")
}
