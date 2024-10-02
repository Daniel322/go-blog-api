package diff_utils

import (
	"reflect"
)

type FindArrResult[T comparable] struct {
	Index  int `json:"index"`
	Result T   `json:"result"`
}

func FindInArr[T comparable](arr []T, key string, findParam any) FindArrResult[T] {
	for i := 0; i < len(arr); i++ {
		checkParam := getField(arr[i], key)
		if checkParam == findParam {
			return FindArrResult[T]{Index: i, Result: arr[i]}
		}
	}

	return FindArrResult[T]{}
}

func getField[T comparable](v T, field string) reflect.Value {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r)
	z := f.FieldByName(field)

	return z
}
