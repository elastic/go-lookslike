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

package llpath

import (
	"reflect"
	"testing"
)

func TestPathComponentType_String(t *testing.T) {
	tests := []struct {
		name string
		pct  PathComponentType
		want string
	}{
		{
			"Should return the correct type",
			pcMapKey,
			"map",
		},
		{
			"Should return the correct type",
			pcSliceIdx,
			"slice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pct.String(); got != tt.want {
				t.Errorf("PathComponentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathComponent_String(t *testing.T) {
	type fields struct {
		Type  PathComponentType
		Key   string
		Index int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"map[string]interface{} key should return a literal",
			fields{pcMapKey, "foo", 0},
			"foo",
		},
		{
			"Array index should return a bracketed number",
			fields{pcSliceIdx, "", 123},
			"[123]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := PathComponent{
				Type:  tt.fields.Type,
				Key:   tt.fields.Key,
				Index: tt.fields.Index,
			}
			if got := pc.String(); got != tt.want {
				t.Errorf("PathComponent.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_ExtendSlice(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		p    Path
		args args
		want Path
	}{
		{
			"Extending an empty slice",
			Path{},
			args{123},
			Path{PathComponent{pcSliceIdx, "", 123}},
		},
		{
			"Extending a non-empty slice",
			Path{PathComponent{pcMapKey, "foo", -1}},
			args{123},
			Path{PathComponent{pcMapKey, "foo", -1}, PathComponent{pcSliceIdx, "", 123}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ExtendSlice(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Path.ExtendSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_ExtendMap(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		p    Path
		args args
		want Path
	}{
		{
			"Extending an empty slice",
			Path{},
			args{"foo"},
			Path{PathComponent{pcMapKey, "foo", -1}},
		},
		{
			"Extending a non-empty slice",
			Path{}.ExtendMap("foo"),
			args{"bar"},
			Path{PathComponent{pcMapKey, "foo", -1}, PathComponent{pcMapKey, "bar", -1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ExtendMap(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Path.ExtendMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Concat(t *testing.T) {
	tests := []struct {
		name string
		p    Path
		arg  Path
		want Path
	}{
		{
			"simple",
			Path{}.ExtendMap("foo"),
			Path{}.ExtendSlice(123),
			Path{}.ExtendMap("foo").ExtendSlice(123),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Concat(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Path.Concat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_String(t *testing.T) {
	tests := []struct {
		name string
		p    Path
		want string
	}{
		{
			"empty",
			Path{},
			"",
		},
		{
			"one element",
			Path{}.ExtendMap("foo"),
			"foo",
		},
		{
			"complex",
			Path{}.ExtendMap("foo").ExtendSlice(123).ExtendMap("bar"),
			"foo.[123].bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Path.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Last(t *testing.T) {
	tests := []struct {
		name string
		p    Path
		want *PathComponent
	}{
		{
			"empty Path",
			Path{},
			nil,
		},
		{
			"one element",
			Path{}.ExtendMap("foo"),
			&PathComponent{pcMapKey, "foo", -1},
		},
		{
			"many elements",
			Path{}.ExtendMap("foo").ExtendMap("bar").ExtendSlice(123),
			&PathComponent{pcSliceIdx, "", 123},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Path.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_GetFrom(t *testing.T) {
	fooPath := Path{}.ExtendMap("foo")
	complexPath := Path{}.ExtendMap("foo").ExtendSlice(0).ExtendMap("bar").ExtendSlice(1)
	tests := []struct {
		name       string
		p          Path
		arg        map[string]interface{}
		wantValue  interface{}
		wantExists bool
	}{
		{
			"simple present",
			fooPath,
			map[string]interface{}{"foo": "bar"},
			"bar",
			true,
		},
		{
			"simple missing",
			fooPath,
			map[string]interface{}{},
			nil,
			false,
		},
		{
			"complex present",
			complexPath,
			map[string]interface{}{"foo": []interface{}{map[string]interface{}{"bar": []string{"bad", "good"}}}},
			"good",
			true,
		},
		{
			"complex missing",
			complexPath,
			map[string]interface{}{},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotExists := tt.p.GetFrom(reflect.ValueOf(tt.arg))
			if tt.wantExists {
				if !reflect.DeepEqual(gotValue.Interface(), tt.wantValue) {
					t.Errorf("Path.GetFrom() gotValue = %v, want %v", gotValue, tt.wantValue)
				}
			} else {
				zero := reflect.Value{}
				if !(gotValue == zero) {
					t.Errorf("Path.GetFrom() gotValue = %v, expected zero value", gotValue)
				}
			}
			if gotExists != tt.wantExists {
				t.Errorf("Path.GetFrom() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		wantP   Path
		wantErr bool
	}{
		{
			"simple",
			"foo",
			Path{}.ExtendMap("foo"),
			false,
		},
		{
			"complex",
			"foo.[0].bar.[1].baz",
			Path{}.ExtendMap("foo").ExtendSlice(0).ExtendMap("bar").ExtendSlice(1).ExtendMap("baz"),
			false,
		},
		{
			"empty string",
			"",
			Path{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, err := ParsePath(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("ParsePath() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}
