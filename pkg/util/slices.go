package util

import (
	"fmt"
	"reflect"
	"strings"
)

func GetIndexWithFieldValue[T any](arr []T, fieldName string, targetValue interface{}) (int, error) {
	for i, elem := range arr {
		elemValue := reflect.ValueOf(elem)
		fieldValue := elemValue.FieldByName(fieldName)
		if fieldValue.IsValid() && (fieldValue.Type().AssignableTo(reflect.TypeOf(targetValue)) || reflect.TypeOf(targetValue).AssignableTo(fieldValue.Type())) && fieldValue.Interface() == targetValue {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Target value not found in slice")
}

func RemoveAtIndex[T any](slice []T, s int) ([]T, error) {
	if slice == nil {
		return nil, fmt.Errorf("input slice is nil")
	}
	if s < 0 || s >= len(slice) {
		return nil, fmt.Errorf("index out of range")
	}
	return append(slice[:s], slice[s+1:]...), nil
}

func SliceHas[T comparable](slice []T, x T) bool {
	if slice == nil {
		return false
	}

	for _, element := range slice {
		if element == x {
			return true
		}
	}

	return false
}

func FilterSliceWithOther[T comparable](src []T, filterWith []T) []T {
	m := map[T]struct{}{}
	for _, v := range filterWith {
		m[v] = struct{}{}
	}

	filtered := []T{}
	for _, v := range src {
		if _, ok := m[v]; !ok {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

func JoinIfNotEmpty(slice []string, sep string) string {
	nonEmptySlice := []string{}

	for _, val := range slice {
		if val == "" {
			continue
		}

		nonEmptySlice = append(nonEmptySlice, val)
	}

	return strings.Join(nonEmptySlice, sep)
}

func FirstNotNil(slice ...any) any {
	for _, val := range slice {
		if val != nil {
			return val
		}
	}

	return nil
}

func StringContainsAny(slice []string, x string) bool {
	for _, element := range slice {
		if strings.Contains(x, element) {
			return true
		}
	}

	return false
}

func ArrayStringContainsIgnoreCase(arr []string, str string) bool {
	for _, s := range arr {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

/*
Ex:
var a int
var b []int
a = 1
b = []int{1, 2, 3}
InArray(a, b) // return true, 0
intinya returnya bool sama index nya jika ditemukan
*/
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
