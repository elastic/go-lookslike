package isdefs

import (
	"fmt"
	"reflect"

	"github.com/elastic/lookslike/lookslike/paths"
	"github.com/elastic/lookslike/lookslike/results"
	"github.com/elastic/lookslike/lookslike/util"
	"github.com/elastic/lookslike/lookslike/validator"
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

// IsSliceOf validates that the array at the given key is an array of objects all validatable
// via the given validator.Validator.
func IsSliceOf(validator validator.Validator) IsDef {
	return Is("slice", func(path paths.Path, v interface{}) *results.Results {
		if reflect.TypeOf(v).Kind() != reflect.Slice {
			return results.SimpleResult(path, false, "Expected slice at given path")
		}
		vSlice := util.InterfaceToSliceOfInterfaces(v)

		res := results.NewResults()

		for idx, curV := range vSlice {
			var validatorRes *results.Results
			validatorRes = validator(curV)
			res.MergeUnderPrefix(path.ExtendSlice(idx), validatorRes)
		}

		return res
	})
}

// IsAny takes a variable number of IsDef's and combines them with a logical OR. If any single definition
// matches the key will be marked as valid.
func IsAny(of ...IsDef) IsDef {
	names := make([]string, len(of))
	for i, def := range of {
		names[i] = def.Name
	}
	isName := fmt.Sprintf("either %#v", names)

	return Is(isName, func(path paths.Path, v interface{}) *results.Results {
		for _, def := range of {
			vr := def.Check(path, v, true)
			if vr.Valid {
				return vr
			}
		}

		return results.SimpleResult(
			path,
			false,
			fmt.Sprintf("Value was none of %#v, actual value was %#v", names, v),
		)
	})
}

// IsUnique instances are used in multiple spots, flagging a value as being in error if it's seen across invocations.
// To use it, assign IsUnique to a variable, then use that variable multiple times in a validator.Map.
func IsUnique() IsDef {
	return ScopedIsUnique().IsUniqueTo("")
}

// UniqScopeTracker is represents the tracking data for invoking IsUniqueTo.
type UniqScopeTracker map[interface{}]string

// IsUniqueTo validates that the given value is only ever seen within a single namespace.
func (ust UniqScopeTracker) IsUniqueTo(namespace string) IsDef {
	return Is("unique", func(path paths.Path, v interface{}) *results.Results {
		for trackerK, trackerNs := range ust {
			hasNamespace := len(namespace) > 0
			if reflect.DeepEqual(trackerK, v) && (!hasNamespace || namespace != trackerNs) {
				return results.SimpleResult(path, false, "Value '%v' is repeated", v)
			}
		}

		ust[v] = namespace
		return results.ValidResult(path)
	})
}

// ScopedIsUnique returns a new scope for uniqueness checks.
func ScopedIsUnique() UniqScopeTracker {
	return UniqScopeTracker{}
}
