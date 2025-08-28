package http

import (
	"case-management/internal/app/usecase"
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
	Queue      *QueueHandler
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
	QueueUC      *usecase.QueueUsecase
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
		Queue:      &QueueHandler{UserCase: *deps.QueueUC},
	}
}
