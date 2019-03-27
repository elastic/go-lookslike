package isdefs

import (
	"time"

	"github.com/elastic/lookslike/lookslike/paths"
	"github.com/elastic/lookslike/lookslike/results"
)

// IsEqualToTime ensures that the actual value is the given time, regardless of zone.
func IsEqualToTime(to time.Time) IsDef {
	return Is("equal to time", func(path paths.Path, v interface{}) *results.Results {
		actualTime, ok := v.(time.Time)
		if !ok {
			return results.SimpleResult(path, false, "Value %t was not a time.Time", v)
		}

		if actualTime.Equal(to) {
			return results.ValidResult(path)
		}

		return results.SimpleResult(path, false, "actual(%v) != expected(%v)", actualTime, to)
	})
}
