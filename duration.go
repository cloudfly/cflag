package confo

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Duration is a flag for holding duration.
type Duration struct {
	// Msecs contains parsed duration in milliseconds.
	Msecs       int64
	valueString string
}

// String implements flag.Value interface
func (d *Duration) String() string {
	return d.valueString
}

// Set implements confo.Value interface
func (d *Duration) Set(value string) error {
	// An attempt to parse value in months.
	minutes, err := strconv.ParseFloat(value, 64)
	if err == nil {
		if minutes < 0 {
			return fmt.Errorf("duration months cannot be negative; got %g", minutes)
		}
		d.Msecs = int64(minutes * msecsPerMinute)
		d.valueString = value
		return nil
	}
	f, err := strconv.ParseFloat(value[:len(value)-1], 64)
	if err != nil {
		return fmt.Errorf("invalid duration value '%s'", value)
	}

	var msecs float64
	// Parse duration.
	switch {
	case strings.HasSuffix(value, "s"):
		msecs = f * msecsPerSecond
	case strings.HasSuffix(value, "m"):
		msecs = f * msecsPerMinute
	case strings.HasSuffix(value, "h"):
		msecs = f * msecsPerHour
	case strings.HasSuffix(value, "d"):
		msecs = f * msecsPerDay
	case strings.HasSuffix(value, "M"):
		msecs = f * msecsPerMonth
	case strings.HasSuffix(value, "y"):
		msecs = f * msecsPerYear
	default:
		return fmt.Errorf("can not parse duration value '%s': unknown duration unit", value)
	}
	d.Msecs = int64(msecs)
	d.valueString = value
	return nil
}

func (d *Duration) UnmarshalJSON(value []byte) error {
	return d.Set(string(value))
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	return d.Set(node.Value)
}

const msecsPerSecond = float64(1000)
const msecsPerMinute = 60 * msecsPerSecond
const msecsPerHour = 60 * msecsPerMinute
const msecsPerDay = 24 * msecsPerHour
const msecsPerMonth = 31 * msecsPerDay
const msecsPerYear = 365 * msecsPerDay
