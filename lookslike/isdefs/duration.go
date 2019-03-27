package isdefs

import (
	"fmt"
	"time"

	"github.com/elastic/lookslike/lookslike/paths"
	"github.com/elastic/lookslike/lookslike/results"
)

// IsDuration tests that the given value is a duration.
var IsDuration = Is("is a duration", func(path paths.Path, v interface{}) *results.Results {
	if _, ok := v.(time.Duration); ok {
		return results.ValidResult(path)
	}
	return results.SimpleResult(
		path,
		false,
		fmt.Sprintf("Expected a time.duration, got '%v' which is a %T", v, v),
	)
})
