package isdefs

import (
	"testing"
	"time"
)

func TestIsDuration(t *testing.T) {
	id := IsDuration

	assertIsDefValid(t, id, time.Duration(1))
	assertIsDefInvalid(t, id, "foo")
}
