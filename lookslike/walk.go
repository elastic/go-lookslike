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
	"github.com/elastic/lookslike/lookslike/paths"
	"github.com/elastic/lookslike/lookslike/util"
	"github.com/elastic/lookslike/lookslike/validator"
	"reflect"
)

type walkObserverInfo struct {
	key     paths.PathComponent
	value   interface{}
	root    validator.Map
	path    paths.Path
}

// walkObserver functions run once per object in the tree.
type walkObserver func(info walkObserverInfo) error

// walk determine if in is a `validator.Map` or a `Slice` and traverse it if so, otherwise will
// treat it as a scalar and invoke the walk observer on the input value directly.
func walk(in interface{}, expandPaths bool, wo walkObserver) error {
	switch in.(type) {
	case validator.Map:
		return walkMap(in.(validator.Map), expandPaths, wo)
	case validator.Slice:
		return walkSlice(in.(validator.Slice), expandPaths, wo)
	case []interface{}:
		return walkSlice(validator.Slice(in.([]interface{})), expandPaths, wo)
	default:
		return walkScalar(in.(validator.Scalar), expandPaths, wo)
	}
}

// walkvalidator.Map is a shorthand way to walk a tree with a map as the root.
func walkMap(m validator.Map, expandPaths bool, wo walkObserver) error {
	return walkFullMap(m, m, paths.Path{}, expandPaths, wo)
}

// walkSlice walks the provided root slice.
func walkSlice(s validator.Slice, expandPaths bool, wo walkObserver) error {
	return walkFullSlice(s, validator.Map{}, paths.Path{}, expandPaths, wo)
}

func walkScalar(s validator.Scalar, expandPaths bool, wo walkObserver) error {
	return wo(walkObserverInfo{
		value: s,
		key:   paths.PathComponent{},
		root:  validator.Map{},
		path:  paths.Path{},
	})
}

func walkFull(o interface{}, root validator.Map, path paths.Path, expandPaths bool, wo walkObserver) (err error) {
	lastPathComponent := path.Last()
	if lastPathComponent == nil {
		// In the case of a slice we can have an empty path
		if _, ok := o.(validator.Slice); ok {
			lastPathComponent = &paths.PathComponent{}
		} else {
			panic("Attempted to traverse an empty Path on a validator.Map in lookslike.walkFull, this should never happen.")
		}
	}

	err = wo(walkObserverInfo{*lastPathComponent, o, root, path})
	if err != nil {
		return err
	}

	switch reflect.TypeOf(o).Kind() {
	case reflect.Map:
		converted := util.InterfaceToMap(o)
		err := walkFullMap(converted, root, path, expandPaths, wo)
		if err != nil {
			return err
		}
	case reflect.Slice:
		converted := util.SliceToSliceOfInterfaces(o)

		for idx, v := range converted {
			newPath := path.ExtendSlice(idx)
			err := walkFull(v, root, newPath, expandPaths, wo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// walkFull walks the given validator.Map tree.
func walkFullMap(m validator.Map, root validator.Map, p paths.Path, expandPaths bool, wo walkObserver) (err error) {
	for k, v := range m {
		var newPath paths.Path
		if !expandPaths {
			newPath = p.ExtendMap(k)
		} else {
			additionalPath, err := paths.ParsePath(k)
			if err != nil {
				return err
			}
			newPath = p.Concat(additionalPath)
		}

		err = walkFull(v, root, newPath, expandPaths, wo)
		if err != nil {
			return err
		}
	}

	return nil
}

func walkFullSlice(s validator.Slice, root validator.Map, p paths.Path, expandPaths bool, wo walkObserver) (err error) {
	for idx, v := range s {
		var newPath paths.Path
		newPath = p.ExtendSlice(idx)

		err = walkFull(v, root, newPath, expandPaths, wo)
		if err != nil {
			return err
		}
	}

	return nil
}
