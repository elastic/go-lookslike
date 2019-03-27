package isdefs

import (
	"fmt"
	"github.com/elastic/lookslike/lookslike/paths"
	"github.com/elastic/lookslike/lookslike/results"
	"regexp"
	"strings"
)

// isStrCheck is a helper for IsDefs that must assert that the value is a string first.
func isStrCheck(path paths.Path, v interface{}) (str string, errorResults *results.Results) {
	strV, ok := v.(string)

	if !ok {
		return "", results.SimpleResult(
			path,
			false,
			fmt.Sprintf("Unable to convert '%v' to string", v),
		)
	}

	return strV, nil
}

// IsString checks that the given value is a string.
var IsString = Is("is a string", func(path paths.Path, v interface{}) *results.Results {
	_, errorResults := isStrCheck(path, v)
	if errorResults != nil {
		return errorResults
	}

	return results.ValidResult(path)
})

// IsNonEmptyString checks that the given value is a string and has a length > 1.
var IsNonEmptyString = Is("is a non-empty string", func(path paths.Path, v interface{}) *results.Results {
	strV, errorResults := isStrCheck(path, v)
	if errorResults != nil {
		return errorResults
	}

	if len(strV) == 0 {
		return results.SimpleResult(path, false, "String '%s' should not be empty", strV)
	}

	return results.ValidResult(path)
})

// IsStringMatching checks whether a value matches the given regexp.
func IsStringMatching(regexp *regexp.Regexp) IsDef {
	return Is("is string matching regexp", func(path paths.Path, v interface{}) *results.Results {
		strV, errorResults := isStrCheck(path, v)
		if errorResults != nil {
			return errorResults
		}

		if !regexp.MatchString(strV) {
			return results.SimpleResult(
				path,
				false,
				fmt.Sprintf("String '%s' did not match regexp %s", strV, regexp.String()),
			)
		}

		return results.ValidResult(path)
	})
}

// IsStringContaining validates that the the actual value contains the specified substring.
func IsStringContaining(needle string) IsDef {
	return Is("is string containing", func(path paths.Path, v interface{}) *results.Results {
		strV, errorResults := isStrCheck(path, v)
		if errorResults != nil {
			return errorResults
		}

		if !strings.Contains(strV, needle) {
			return results.SimpleResult(
				path,
				false,
				fmt.Sprintf("String '%s' did not contain substring '%s'", strV, needle),
			)
		}

		return results.ValidResult(path)
	})
}

