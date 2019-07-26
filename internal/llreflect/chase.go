package llreflect

import (
	"fmt"
	"reflect"
)

func ChaseValue(v reflect.Value) reflect.Value {
	for (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
		v = v.Elem()
	}
	s := v.String()
	fmt.Sprintf(s)
	return v
}