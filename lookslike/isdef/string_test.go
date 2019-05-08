package isdef

import (
	"regexp"
	"testing"
)

func TestIsString(t *testing.T) {
	assertIsDefValid(t, IsString, "abc")
	assertIsDefValid(t, IsString, "a")
	assertIsDefInvalid(t, IsString, 123)
}

func TestIsNonEmptyString(t *testing.T) {
	assertIsDefValid(t, IsNonEmptyString, "abc")
	assertIsDefValid(t, IsNonEmptyString, "a")
	assertIsDefInvalid(t, IsNonEmptyString, "")
	assertIsDefInvalid(t, IsString, 123)
}

func TestIsStringMatching(t *testing.T) {
	id := IsStringMatching(regexp.MustCompile(`^f`))

	assertIsDefValid(t, id, "fall")
	assertIsDefInvalid(t, id, "potato")
	assertIsDefInvalid(t, IsString, 123)
}

func TestIsStringContaining(t *testing.T) {
	id := IsStringContaining("foo")

	assertIsDefValid(t, id, "foo")
	assertIsDefValid(t, id, "a foo b")
	assertIsDefInvalid(t, id, "a bar b")
	assertIsDefInvalid(t, IsString, 123)
}
