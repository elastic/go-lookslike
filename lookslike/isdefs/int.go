package isdefs

import (
	"fmt"

	"github.com/elastic/lookslike/lookslike/paths"
	"github.com/elastic/lookslike/lookslike/results"
)

func intGtChecker(than int) ValueValidator {
	return func(path paths.Path, v interface{}) *results.Results {
		n, ok := v.(int)
		if !ok {
			msg := fmt.Sprintf("%v is a %T, but was expecting an int!", v, v)
			return results.SimpleResult(path, false, msg)
		}

		if n > than {
			return results.ValidResult(path)
		}

		return results.SimpleResult(
			path,
			false,
			fmt.Sprintf("%v is not greater than %v", n, than),
		)
	}
}

// IsIntGt tests that a value is an int greater than.
func IsIntGt(than int) IsDef {
	return Is("greater than", intGtChecker(than))
}
