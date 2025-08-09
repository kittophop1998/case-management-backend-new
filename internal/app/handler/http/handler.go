package http

import (
	"case-management/internal/app/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User       *UserHandler
	MasterData *MasterDataHandler
	Auth       *AuthHandler
	Permission *PermissionHandler
	Log        *LogHandler
}

var H *Handlers

func InitHandlers(
	userUC *usecase.UserUseCase,
	masterDataUC *usecase.MasterDataUseCase,
	authUC *usecase.AuthUseCase,
	permissionUC *usecase.PermissionUseCase,
	logUC *usecase.LogUseCase,
) {
	H = &Handlers{
		Auth:       &AuthHandler{UseCase: *authUC},
		User:       &UserHandler{UseCase: *userUC},
		MasterData: &MasterDataHandler{UseCase: *masterDataUC},
		Permission: &PermissionHandler{UseCase: *permissionUC},
		Log:        &LogHandler{UseCase: *logUC},
	}
}

func getLimit(ctx *gin.Context) (int, error) {
	limit := 10
	var err error

	limitParam := ctx.Query("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return 0, err
		}
	}

	return limit, nil
}

func getPage(ctx *gin.Context) (int, error) {
	page := 1
	var err error

	pageParam := ctx.Query("page")
	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			return 0, err
		}
	}

	return page, nil
}
