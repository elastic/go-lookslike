package isdef

import (
	"testing"
	"time"

	"github.com/elastic/lookslike/lookslike/llpath"
	"github.com/elastic/lookslike/lookslike/llresult"
	"github.com/stretchr/testify/assert"
)

func assertIsDefValid(t *testing.T, id IsDef, value interface{}) *llresult.Results {
	res := id.Check(llpath.MustParsePath("p"), value, true)

	if !res.Valid {
		assert.Fail(
			t,
			"Expected Valid IsDef",
			"Isdef %#v was not valid for value %#v with error: ", id, value, res.Errors(),
		)
	}
	return res
}

func assertIsDefInvalid(t *testing.T, id IsDef, value interface{}) *llresult.Results {
	res := id.Check(llpath.MustParsePath("p"), value, true)

	if res.Valid {
		assert.Fail(
			t,
			"Expected invalid IsDef",
			"Isdef %#v was should not have been valid for value %#v",
			id,
			value,
		)
	}
	return res
}

func TestIsAny(t *testing.T) {
	id := IsAny(IsEqual("foo"), IsEqual("bar"))

	assertIsDefValid(t, id, "foo")
	assertIsDefValid(t, id, "bar")
	assertIsDefInvalid(t, id, "basta")
}

func TestIsEqual(t *testing.T) {
	id := IsEqual("foo")

	assertIsDefValid(t, id, "foo")
	assertIsDefInvalid(t, id, "bar")
}

func TestRegisteredIsEqual(t *testing.T) {
	// Time equality comes from a registered function
	// so this is a quick way to test registered functions
	now := time.Now()
	id := IsEqual(now)

	assertIsDefValid(t, id, now)
	assertIsDefInvalid(t, id, now.Add(100))
}

func TestIsNil(t *testing.T) {
	assertIsDefValid(t, IsNil, nil)
	assertIsDefInvalid(t, IsNil, "foo")
}
