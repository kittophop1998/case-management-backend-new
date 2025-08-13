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
	Case       *CaseHandler
	Customer   *CustomerHandler
	Dashboard  *DashboardHandler
}

type HandlerDeps struct {
	UserUC       *usecase.UserUseCase
	MasterDataUC *usecase.MasterDataUseCase
	AuthUC       *usecase.AuthUseCase
	PermissionUC *usecase.PermissionUseCase
	LogUC        *usecase.LogUseCase
	CaseUC       *usecase.CaseUseCase
	CustomerUC   *usecase.CustomerUseCase
	DashboardUC  *usecase.DashboardUseCase
}

var H *Handlers

func InitHandlers(deps HandlerDeps) {
	H = &Handlers{
		Auth:       &AuthHandler{UseCase: *deps.AuthUC, UserUseCase: *deps.UserUC},
		User:       &UserHandler{UseCase: *deps.UserUC},
		MasterData: &MasterDataHandler{UseCase: *deps.MasterDataUC},
		Permission: &PermissionHandler{UseCase: *deps.PermissionUC},
		Log:        &LogHandler{UseCase: *deps.LogUC},
		Case:       &CaseHandler{UseCase: *deps.CaseUC},
		Customer:   &CustomerHandler{UseCase: *deps.CustomerUC},
		Dashboard:  &DashboardHandler{UseCase: *deps.DashboardUC},
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
