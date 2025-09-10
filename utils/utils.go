package utils

import (
	"case-management/internal/domain/model"
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//
// ====== String & Bool Helpers ======
//

// NormalizeUserInput ทำให้ username/email เป็น lower-case
func NormalizeUserInput(user *model.CreateUpdateUserRequest) {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
}

// Ptr helpers
func StringPtr(s string) *string { return &s }
func BoolPtr(v bool) *bool       { return &v }

//
// ====== UUID Helpers ======
//

// ParseUUID แปลง string → uuid.UUID
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// ParseOptionalUUID แปลง *string → *uuid.UUID (nil ถ้า string ว่างหรือ nil)
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

// ParseUUIDQueryParam แปลง query param → *uuid.UUID หรือ uuid.UUID (optional/required)
func ParseUUIDQueryParam(c *gin.Context, param string, required bool) (*uuid.UUID, error) {
	val := strings.TrimSpace(c.Query(param))
	if val == "" {
		if required {
			return nil, fmt.Errorf("%s is required", param)
		}
		return nil, nil
	}
	id, err := uuid.Parse(val)
	if err != nil {
		return nil, fmt.Errorf("invalid %s format: %w", param, err)
	}
	return &id, nil
}

// UUIDPtrToStringPtr แปลง *uuid.UUID → *string
func UUIDPtrToStringPtr(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}
	s := u.String()
	return &s
}

// ====== Random Helpers ======
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) (string, error) {
	b := make([]rune, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			return "", err
		}
		b[i] = letterRunes[num.Int64()]
	}
	return string(b), nil
}

// ====== Date Helpers ======
func FormatDate(t *time.Time, layout string) string {
	if t == nil {
		return ""
	}
	return t.Format(layout)
}

func ParseOptionalDate(dateStr *string, layout string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}
	parsed, err := time.Parse(layout, *dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	return &parsed, nil
}

// ====== Reflection Helpers ======
func IsEmpty(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
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
		return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	default:
		return false
	}
}

// ====== Domain Specific ======
func UserNameCenter(user model.User) string {
	return fmt.Sprintf("%s - %s", user.Name, user.Center.Name)
}
