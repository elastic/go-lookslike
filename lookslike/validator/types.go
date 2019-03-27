package validator


// validator.Map is the type used to define schema definitions for Compile and to represent an arbitrary
// map of values of any type.
type Map map[string]interface{}

// Slice is a convenience []interface{} used to declare schema defs. You would typically nest this inside
// a validator.Map as a value, and it would be able to match against any type of non-empty slice.
type Slice []interface{}

// Catchall type for things that aren't assertable to either validator.Map or Slice.
type Scalar interface{}

