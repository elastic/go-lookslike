package validator

import "github.com/elastic/lookslike/lookslike/results"

// Validator is the result of Compile and is run against the map you'd like to test.
type Validator func(interface{}) *results.Results

