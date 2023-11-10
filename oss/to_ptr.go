package oss

import "reflect"

// Ptr returns a pointer to the provided value.
func Ptr[T any](v T) *T {
	return &v
}

// Pta returns a provided value for the pointer.
func Pta(ptr interface{}) interface{} {
	value := reflect.ValueOf(ptr)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return nil
	}
	return value.Elem().Interface()
}

// SliceOfPtrs returns a slice of *T from the specified values.
func SliceOfPtrs[T any](vv ...T) []*T {
	slc := make([]*T, len(vv))
	for i := range vv {
		slc[i] = Ptr(vv[i])
	}
	return slc
}
