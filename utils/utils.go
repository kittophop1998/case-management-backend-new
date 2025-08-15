package utils

import (
	"case-management/internal/domain/model"
	"strings"
)

func NormalizeUserInput(user *model.CreateUpdateUserRequest) {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
}
