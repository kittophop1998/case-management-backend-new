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
// ====== String Helpers ======
//

func NormalizeUserInput(user *model.CreateUpdateUserRequest) {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
}

func StringPtr(s string) *string {
	return &s
}

//
// ====== Bool Helpers ======
//

func Bool(v bool) *bool {
	return &v
}

//
// ====== UUID Helpers ======
//

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

// แปลง query string -> *uuid.UUID
func ParseUUIDParam(c *gin.Context, param string) (*uuid.UUID, error) {
	val := c.Query(param)
	if val == "" {
		return nil, nil
	}
	return ParseOptionalUUID(&val)
}

func UUIDPtrToStringPtr(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}
	str := u.String()
	return &str
}

//
// ====== Random Helpers ======
//

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

//
// ====== Date Helpers ======
//

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

//
// ====== Reflection Helpers ======
//

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

//
// ====== Domain Specific ======
//

func UserNameCenter(user model.User) string {
	return user.Name + " - " + user.Center.Name
}
