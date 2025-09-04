package utils

import (
	"case-management/internal/domain/model"
	"crypto/rand"
	"math/big"
	"reflect"
	"strings"

	"github.com/google/uuid"
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

func RandStringRunes(n int) (string, error) {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		maxBigInt := big.NewInt(int64(len(letterRunes)))
		num, err := rand.Int(rand.Reader, maxBigInt)
		if err != nil {
			return "", err
		}
		b[i] = letterRunes[num.Int64()]
	}
	return string(b), nil
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func ParseOptionalUUID(s *string) (*uuid.UUID, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
