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

package testslike

import (
	"github.com/elastic/lookslike/lookslike"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

// Test takes the output from a Validator invocation and runs test assertions on the result.
// If you are using this library for testing you will probably want to run Test(t, Compile(Map{...}), actual) as a pattern.
func Test(t *testing.T, validator lookslike.Validator, value interface{}) *lookslike.Results {
	r := validator(value)

	if !r.Valid {
		assert.Fail(
			t,
			"lookslike could not validate map",
			"%d errors validating source: \n%s", len(r.Errors()), spew.Sdump(value),
		)
	}

	for _, err := range r.Errors() {
		assert.NoError(t, err)
	}
	return r
}
