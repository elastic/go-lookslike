package isdefs

import (
	"github.com/elastic/lookslike/lookslike/paths"
	"github.com/elastic/lookslike/lookslike/results"
)

// Is creates a named IsDef with the given Checker.
func Is(name string, checker ValueValidator) IsDef {
	return IsDef{Name: name, Checker: checker}
}

// A ValueValidator is used to validate a value in a validator.Map.
type ValueValidator func(path paths.Path, v interface{}) *results.Results

// An IsDef defines the type of Check to do.
// Generally only Name and Checker are set. Optional and CheckKeyMissing are
// needed for weird checks like key presence.
type IsDef struct {
	Name            string
	Checker         ValueValidator
	Optional        bool
	CheckKeyMissing bool
}

// Check runs the IsDef at the given value at the given path
func (id IsDef) Check(path paths.Path, v interface{}, keyExists bool) *results.Results {
	if id.CheckKeyMissing {
		if !keyExists {
			return results.ValidResult(path)
		}

		return results.SimpleResult(path, false, "this key should not exist")
	}

	if !id.Optional && !keyExists {
		return results.KeyMissingResult(path)
	}

	if id.Checker != nil {
		return id.Checker(path, v)
	}

	return results.ValidResult(path)
}

// Optional wraps an IsDef to mark the field's presence as Optional.
func Optional(id IsDef) IsDef {
	id.Name = "Optional " + id.Name
	id.Optional = true
	return id
}
