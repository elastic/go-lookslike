package validator

import "github.com/elastic/lookslike/lookslike/llresult"

// Validator is the result of Compile and is run against the map you'd like to test.
type Validator func(interface{}) *llresult.Results
