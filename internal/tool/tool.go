package tool

import "reflect"

func CheckSliceEqual(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}
