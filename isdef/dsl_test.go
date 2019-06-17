package isdef

import (
	"reflect"
	"testing"

	"github.com/elastic/go-lookslike/llpath"
	"github.com/elastic/go-lookslike/llresult"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsSliceOf(t *testing.T) {
	goodMap := map[string]interface{}{"foo": "bar"}

	isFooBarMap := IsSliceOf(func(i interface{}) *llresult.Results {
		if reflect.DeepEqual(i, goodMap) {
			return llresult.ValidResult(llpath.MustParsePath("foo"))
		}

		return llresult.SimpleResult(llpath.MustParsePath("foo"), false, "did not match")
	})

	goodMapArr := []map[string]interface{}{goodMap, goodMap}

	goodRes := assertIsDefValid(t, isFooBarMap, goodMapArr)
	goodFields := goodRes.Fields
	assert.Len(t, goodFields, 2)
	assert.Contains(t, goodFields, "p.[0].foo")
	assert.Contains(t, goodFields, "p.[1].foo")

	badMap := map[string]interface{}{"foo": "bot"}
	badMapArr := []map[string]interface{}{badMap}

	badRes := assertIsDefInvalid(t, isFooBarMap, badMapArr)
	badFields := badRes.Fields
	assert.Len(t, badFields, 1)
	assert.Contains(t, badFields, "p.[0].foo")
}
func TestIsUnique(t *testing.T) {
	pathFoo := llpath.MustParsePath("foo")
	pathBar := llpath.MustParsePath("bar")

	tests := []struct {
		name    string
		fn      func() *llresult.Results
		isValid bool
	}{
		{
			"IsUnique find dupes",
			func() *llresult.Results {
				u := IsUnique()
				u.Check(pathFoo, "a", true)
				return u.Check(pathBar, "a", true)
			},
			false,
		},
		{
			"IsUnique separate instances don't care about dupes",
			func() *llresult.Results {
				IsUnique().Check(pathFoo, "a", true)
				return IsUnique().Check(pathFoo, "b", true)
			},
			true,
		},
		{
			"IsUniqueTo duplicates across namespaces fail",
			func() *llresult.Results {
				s := ScopedIsUnique()
				s.IsUniqueTo("test").Check(pathFoo, 1, true)
				return s.IsUniqueTo("test2").Check(pathFoo, 1, true)
			},
			false,
		},

		{
			"IsUniqueTo duplicates within a namespace succeeds",
			func() *llresult.Results {
				s := ScopedIsUnique()
				s.IsUniqueTo("test").Check(pathFoo, 1, true)
				return s.IsUniqueTo("test").Check(pathBar, 1, true)
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.isValid {
				assert.True(t, tt.fn().Valid)
			} else {
				require.False(t, tt.fn().Valid)
			}
		})
	}
}
