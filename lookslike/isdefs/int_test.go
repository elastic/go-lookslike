package isdefs

import "testing"

func TestIsIntGt(t *testing.T) {
	id := IsIntGt(100)

	assertIsDefValid(t, id, 101)
	assertIsDefInvalid(t, id, 100)
	assertIsDefInvalid(t, id, 99)
}
