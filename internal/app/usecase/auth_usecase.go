package usecase

import (
	"case-management/infrastructure/auth"
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"case-management/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gopkg.in/ldap.v2"
)

type AuthUseCase struct {
	logUsecase *LogUseCase
	repo       repository.UserRepository
}

func NewAuthUseCase(logUsecase *LogUseCase, repo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		logUsecase: logUsecase,
		repo:       repo,
	}
}

func (a *AuthUseCase) Login(ctx *gin.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := a.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user.UserTypes == "local" && req.Username == user.Username && req.Password == user.Password {
		return a.loginLocal(ctx, user)
	}

	// ----- LDAP Authentication -----
	if err := a.authenticateWithLDAP(req.Username, req.Password); err != nil {
		return nil, err
	}

	accessToken, err := a.GenerateToken(24*time.Hour, &auth.Metadata{
		UserId:     user.ID,
		Name:       user.Name,
		CenterId:   user.Center.ID,
		CenterName: user.Center.Name,
		SectionId:  user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.GenerateToken(3*24*time.Hour, &auth.Metadata{
		UserId:     user.ID,
		Name:       user.Name,
		CenterId:   user.Center.ID,
		CenterName: user.Center.Name,
		SectionId:  user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	logs := &model.AccessLogs{
		UserID:        user.ID,
		Action:        "login",
		Details:       nil,
		CreatedAt:     time.Now(),
		Username:      user.Username,
		LogonDatetime: time.Now(),
		LoginSuccess:  utils.Bool(true),
	}

	if err := a.logUsecase.SaveLoginEvent(ctx, logs); err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthUseCase) authenticateWithLDAP(username, password string) error {
	conn, err := ldap.Dial("tcp", "ldap.example.com:389")
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Bind("HEADOFFICE\\"+username, password); err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) GenerateToken(ttl time.Duration, metadata *auth.Metadata) (string, error) {
	claims := &auth.JwtClaims{
		Metadata: *metadata,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(ttl).Unix(),
			Issuer:    "casemanagement",
		},
	}

	secretKey := []byte("casemanagement_secret_key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (a *AuthUseCase) Logout(ctx *gin.Context) error {
	userIdStr := ctx.GetString("userId")
	username := ctx.GetString("username")

	var userID uuid.UUID
	if id, err := uuid.Parse(userIdStr); err == nil {
		userID = id
	} else {
		return err
	}

	logs := &model.AccessLogs{
		UserID:        userID,
		Action:        "logout",
		Details:       nil,
		CreatedAt:     time.Now(),
		Username:      username,
		LogonDatetime: time.Now(),
		LoginSuccess:  utils.Bool(true),
	}

	if err := a.logUsecase.SaveLoginEvent(ctx, logs); err != nil {
		return err
	}

	return nil
}

// Login for local user
func (a *AuthUseCase) loginLocal(ctx *gin.Context, user *model.User) (*model.LoginResponse, error) {
	accesstoken, err := a.GenerateToken(24*time.Hour, &auth.Metadata{
		UserId:     user.ID,
		Name:       user.Name,
		CenterId:   user.Center.ID,
		CenterName: user.Center.Name,
		SectionId:  user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.GenerateToken(3*24*time.Hour, &auth.Metadata{
		UserId:     user.ID,
		Name:       user.Name,
		CenterId:   user.Center.ID,
		CenterName: user.Center.Name,
		SectionId:  user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	logs := &model.AccessLogs{
		UserID:        user.ID,
		Action:        "login",
		Details:       nil,
		CreatedAt:     time.Now(),
		Username:      user.Username,
		LogonDatetime: time.Now(),
		LoginSuccess:  utils.Bool(true),
	}

	if err := a.logUsecase.SaveLoginEvent(ctx, logs); err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accesstoken,
		RefreshToken: refreshToken,
	}, nil
}
