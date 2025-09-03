package utils

import (
	"case-management/internal/domain/model"
	"reflect"
	"strings"
)

func NormalizeUserInput(user *model.CreateUpdateUserRequest) {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
}

func Bool(v bool) *bool {
	return &v
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	// Dereference pointers & interfaces
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return true
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.String:
		return strings.TrimSpace(val.String()) == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return val.Len() == 0
	case reflect.Struct:
		// เช็ค struct ว่าง (all zero values)
		return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	default:
		return false
	}
}
