package http

import "case-management/internal/app/usecase"

type Handlers struct {
	User       *UserHandler
	MasterData *MasterDataHandler
	Auth       *AuthHandler
}

var H *Handlers

func InitHandlers(
	userUC *usecase.UserUseCase,
	masterDataUC *usecase.MasterDataUseCase,
	authUC *usecase.AuthUseCase,
) {
	H = &Handlers{
		Auth:       &AuthHandler{useCase: *authUC},
		User:       &UserHandler{useCase: *userUC},
		MasterData: &MasterDataHandler{UseCase: *masterDataUC},
	}
}
