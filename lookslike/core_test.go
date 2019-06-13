// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lookslike

import (
	"regexp"
	"testing"
	"time"

	"github.com/elastic/go-lookslike/lookslike/isdef"
	"github.com/elastic/go-lookslike/lookslike/llpath"
	"github.com/elastic/go-lookslike/lookslike/llresult"
	"github.com/elastic/go-lookslike/lookslike/validator"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func assertValidator(t *testing.T, validator validator.Validator, input map[string]interface{}) {
	res := validator(input)
	assertResults(t, res)
}

// assertResults validates the schema passed successfully.
func assertResults(t *testing.T, r *llresult.Results) *llresult.Results {
	for _, err := range r.Errors() {
		assert.NoError(t, err)
	}
	return r
}

func TestFlat(t *testing.T) {
	// Test map[string]interface{} as a user would more likely use
	m := map[string]interface{}{
		"foo": "bar",
		"baz": 1,
	}

	vResults := MustCompile(map[string]interface{}{
		"foo": "bar",
		"baz": isdef.IsIntGt(0),
	})(m)

	assertResults(t, vResults)
}

func TestBadFlat(t *testing.T) {
	m := map[string]interface{}{}

	fakeT := new(testing.T)

	vResults := MustCompile(map[string]interface{}{
		"notafield": isdef.IsDuration,
	})(m)

	assertResults(fakeT, vResults)

	assert.True(t, fakeT.Failed())

	result := vResults.Fields["notafield"][0]
	assert.False(t, result.Valid)
	assert.Equal(t, result, llresult.KeyMissingVR)
}

func TestInterface(t *testing.T) {
	results := MustCompile(isdef.IsEqual(42))(42)
	assertResults(t, results)
}

func TestBadInterface(t *testing.T) {
	fakeT := new(testing.T)

	results := MustCompile(isdef.IsEqual(42))(-1)

	assertResults(fakeT, results)

	assert.True(t, fakeT.Failed())

	result := results.Fields[""][0]
	assert.False(t, result.Valid)
}

func TestInterfaceTypeMismatch(t *testing.T) {
	fakeT := new(testing.T)

	results := MustCompile(isdef.IsEqual(42))("foo")

	assertResults(fakeT, results)

	assert.True(t, fakeT.Failed())

	result := results.Fields[""][0]
	assert.False(t, result.Valid)
}

func TestSlice(t *testing.T) {
	actual := []interface{}{42, time.Second, "admiral akbar"}
	results := MustCompile([]interface{}{42, isdef.IsDuration, isdef.IsStringMatching(regexp.MustCompile("bar"))})(actual)
	assertResults(t, results)
}

func TestBadSlice(t *testing.T) {
	fakeT := new(testing.T)

	actual := []interface{}{42, time.Second, "admiral akbar"}
	results := MustCompile([]interface{}{42, isdef.IsDuration, isdef.IsStringMatching(regexp.MustCompile("NOTHERE"))})(actual)

	assertResults(fakeT, results)

	assert.True(t, fakeT.Failed())

	assert.True(t, results.Fields["[0]"][0].Valid)
	assert.True(t, results.Fields["[1]"][0].Valid)
	assert.False(t, results.Fields["[2]"][0].Valid)
}

func TestSliceStrictness(t *testing.T) {
	// Test that different slice lengths cause a failure
	fakeT := new(testing.T)

	actual := []interface{}{42, time.Second, "admiral akbar", "EXTRA"}
	// One less item than the real thing
	results := MustCompile([]interface{}{42, isdef.IsDuration, isdef.IsStringMatching(regexp.MustCompile("bar"))})(actual)

	assertResults(fakeT, results)

	assert.True(t, fakeT.Failed())

	assert.True(t, results.Fields["[0]"][0].Valid)
	assert.True(t, results.Fields["[1]"][0].Valid)
	assert.True(t, results.Fields["[2]"][0].Valid)
}

func TestPrimitiveSlice(t *testing.T) {
	actual := []int{1, 1, 2, 3}
	results := MustCompile([]interface{}{1, 1, 2, 3})(actual)
	assertResults(t, results)
}

func TestBadPrimitiveSlice(t *testing.T) {
	fakeT := new(testing.T)

	actual := []int{1, 2, 3}
	results := MustCompile([]interface{}{1, 1, 1})(actual)

	assertResults(fakeT, results)

	assert.True(t, fakeT.Failed())

	assert.True(t, results.Fields["[0]"][0].Valid)
	assert.False(t, results.Fields["[1]"][0].Valid)
	assert.False(t, results.Fields["[2]"][0].Valid)
}

func TestNested(t *testing.T) {
	m := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
			"dur": time.Duration(100),
		},
	}

	results := MustCompile(map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
		"foo.dur": isdef.IsDuration,
	})(m)

	assertResults(t, results)

	assert.Len(t, results.Fields, 2, "One result per matcher")
}

func TestComposition(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"baz": "bot",
	}

	fooValidator := MustCompile(map[string]interface{}{"foo": "bar"})
	bazValidator := MustCompile(map[string]interface{}{"baz": "bot"})
	composed := Compose(fooValidator, bazValidator)

	// Test that the validators work individually
	assertValidator(t, fooValidator, m)
	assertValidator(t, bazValidator, m)

	// Test that the composition of them works
	assertValidator(t, composed, m)

	composedRes := composed(m)
	assert.Len(t, composedRes.Fields, 2)

	badValidator := MustCompile(map[string]interface{}{"notakey": "blah"})
	badComposed := Compose(badValidator, composed)

	fakeT := new(testing.T)
	assertValidator(fakeT, badComposed, m)
	badComposedRes := badComposed(m)

	assert.Len(t, badComposedRes.Fields, 3)
	assert.True(t, fakeT.Failed())
}

func TestStrictFunc(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"baz": "bot",
		"nest": map[string]interface{}{
			"very": map[string]interface{}{
				"deep": "true",
			},
		},
	}

	validValidator := MustCompile(map[string]interface{}{
		"foo": "bar",
		"baz": "bot",
		"nest": map[string]interface{}{
			"very": map[string]interface{}{
				"deep": "true",
			},
		},
	})

	assertValidator(t, validValidator, m)

	partialValidator := MustCompile(map[string]interface{}{
		"foo": "bar",
	})

	// Should pass, since this is not a strict Check
	assertValidator(t, partialValidator, m)

	res := Strict(partialValidator)(m)

	assert.Equal(t, []llresult.ValueResult{llresult.StrictFailureVR}, res.DetailedErrors().Fields["baz"])
	assert.Equal(t, []llresult.ValueResult{llresult.StrictFailureVR}, res.DetailedErrors().Fields["nest.very.deep"])
	assert.Nil(t, res.DetailedErrors().Fields["bar"])
	assert.False(t, res.Valid)
}

func TestExistence(t *testing.T) {
	m := map[string]interface{}{
		"exists": "foo",
	}

	validator := MustCompile(map[string]interface{}{
		"exists": isdef.KeyPresent,
		"non":    isdef.KeyMissing,
	})

	assertValidator(t, validator, m)
}

func TestComplex(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"hash": map[string]interface{}{
			"baz": 1,
			"bot": 2,
			"deep_hash": map[string]interface{}{
				"qux": "quark",
			},
		},
		"slice": []string{"pizza", "pasta", "and more"},
		"empty": nil,
		"arr":   []map[string]interface{}{{"foo": "bar"}, {"foo": "baz"}},
	}

	validator := MustCompile(map[string]interface{}{
		"foo": "bar",
		"hash": map[string]interface{}{
			"baz": 1,
			"bot": 2,
			"deep_hash": map[string]interface{}{
				"qux": "quark",
			},
		},
		"slice":        []string{"pizza", "pasta", "and more"},
		"empty":        isdef.KeyPresent,
		"doesNotExist": isdef.KeyMissing,
		"arr":          isdef.IsSliceOf(MustCompile(map[string]interface{}{"foo": isdef.IsStringContaining("a")})),
	})

	assertValidator(t, validator, m)
}

func TestLiteralArray(t *testing.T) {
	m := map[string]interface{}{
		"a": []interface{}{
			[]interface{}{1, 2, 3},
			[]interface{}{"foo", "bar"},
			"hello",
		},
	}

	validator := MustCompile(map[string]interface{}{
		"a": []interface{}{
			[]interface{}{1, 2, 3},
			[]interface{}{"foo", "bar"},
			"hello",
		},
	})

	goodRes := validator(m)

	assertResults(t, goodRes)
	// We evaluate multidimensional slice as a single field for now
	// This is kind of easier, but maybe we should do our own traversal later.
	assert.Len(t, goodRes.Fields, 6)
}

func TestStringSlice(t *testing.T) {
	m := map[string]interface{}{
		"a": []string{"a", "b"},
	}

	validator := MustCompile(map[string]interface{}{
		"a": []string{"a", "b"},
	})

	goodRes := validator(m)

	assertResults(t, goodRes)
	// We evaluate multidimensional slices as a single field for now
	// This is kind of easier, but maybe we should do our own traversal later.
	assert.Len(t, goodRes.Fields, 2)
}

func TestEmptySlice(t *testing.T) {
	// In the case of an empty Slice, the validator will compare slice type
	// In this case we're treating the slice as a value and doing a literal comparison
	// Users should use an IsDef testing for an empty slice (that can use reflection)
	// if they need something else.
	m := map[string]interface{}{
		"a": []interface{}{},
		"b": []string{},
	}

	validator := MustCompile(map[string]interface{}{
		"a": []interface{}{},
		"b": []string{},
	})

	goodRes := validator(m)

	assertResults(t, goodRes)
	assert.Len(t, goodRes.Fields, 2)
}

func TestLiteralMdSlice(t *testing.T) {
	m := map[string]interface{}{
		"a": [][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
	}

	v := MustCompile(map[string]interface{}{
		"a": [][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
	})

	goodRes := v(m)

	assertResults(t, goodRes)
	// We evaluate multidimensional slices as a single field for now
	// This is kind of easier, but maybe we should do our own traversal later.
	assert.Len(t, goodRes.Fields, 6)

	badValidator := Strict(MustCompile(map[string]interface{}{
		"a": [][]int{
			{1, 2, 3},
		},
	}))

	badRes := badValidator(m)

	assert.False(t, badRes.Valid)
	assert.Len(t, badRes.Fields, 7)
	// The reason the len is 4 is that there is 1 extra slice + 4 values.
	assert.Len(t, badRes.Errors(), 4)
}

func TestSliceOfIsDefs(t *testing.T) {
	m := map[string]interface{}{
		"a": []int{1, 2, 3},
		"b": []interface{}{"foo", "bar", 3},
	}

	goodV := MustCompile(map[string]interface{}{
		"a": []interface{}{isdef.IsIntGt(0), isdef.IsIntGt(1), 3},
		"b": []interface{}{isdef.IsStringContaining("o"), "bar", isdef.IsIntGt(2)},
	})

	assertValidator(t, goodV, m)

	badV := MustCompile(map[string]interface{}{
		"a": []interface{}{isdef.IsIntGt(100), isdef.IsIntGt(1), 3},
		"b": []interface{}{isdef.IsStringContaining("X"), "bar", isdef.IsIntGt(2)},
	})
	badRes := badV(m)

	assert.False(t, badRes.Valid)
	assert.Len(t, badRes.Errors(), 2)
}

func TestMatchArrayAsValue(t *testing.T) {
	m := map[string]interface{}{
		"a": []int{1, 2, 3},
		"b": []interface{}{"foo", "bar", 3},
	}

	goodV := MustCompile(map[string]interface{}{
		"a": []int{1, 2, 3},
		"b": []interface{}{"foo", "bar", 3},
	})

	assertValidator(t, goodV, m)

	badV := MustCompile(map[string]interface{}{
		"a": "robot",
		"b": []interface{}{"foo", "bar", 3},
	})

	badRes := badV(m)

	assert.False(t, badRes.Valid)
	assert.False(t, badRes.Fields["a"][0].Valid)
	for _, f := range badRes.Fields["b"] {
		assert.True(t, f.Valid)
	}
}

func TestInvalidPathIsdef(t *testing.T) {
	badPath := "foo...bar"
	_, err := compile(map[string]interface{}{
		badPath: "invalid",
	})

	assert.Equal(t, llpath.InvalidPathString(badPath), err)
}

// This test is here, not in isdefs because it really is testing core functionality
func TestOptional(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
	}

	validator := MustCompile(map[string]interface{}{
		"non": isdef.Optional(isdef.IsEqual("foo")),
	})

	require.True(t, validator(m).Valid)
}
