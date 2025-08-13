package usecase

import (
	"case-management/infrastructure/auth"
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gopkg.in/ldap.v2"
)

type AuthUseCase struct {
	repo     repository.UserRepository
	authRepo repository.AuthRepository
}

func NewAuthUseCase(repo repository.UserRepository, authRepo repository.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		repo:     repo,
		authRepo: authRepo,
	}
}

func (a *AuthUseCase) Login(ctx *gin.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	if req.Username == "admin" && req.Password == "admin" {
		return a.loginAsAdmin(ctx, req.Username)
	}

	// Authenticate via LDAP
	if err := a.authenticateWithLDAP(req.Username, req.Password); err != nil {
		return nil, err
	}

	// Fetch user from DB
	user, err := a.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.GenerateToken(24*time.Hour, &auth.Metadata{
		UserId:    user.ID,
		Name:      user.Name,
		CenterId:  user.Center.ID,
		SectionId: user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.GenerateToken(3*24*time.Hour, &auth.Metadata{
		UserId:    user.ID,
		Name:      user.Name,
		CenterId:  user.Center.ID,
		SectionId: user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// authenticateWithLDAP authenticates via external LDAP
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

// Token functions (delegate to repository)
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

func (a *AuthUseCase) SaveAccessLog(ctx *gin.Context, username string, success bool) error {

	userIdStr := ctx.GetString("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	return a.authRepo.SaveAccessLog(ctx, &model.AccessLogs{
		UserID:        userId,
		Action:        "login",
		IPAddress:     ctx.ClientIP(),
		UserAgent:     ctx.GetHeader("User-Agent"),
		Details:       nil,
		CreatedAt:     time.Now(),
		Username:      username,
		LogonDatetime: time.Now(),
		LogonResult:   "success",
	})
}

func (a *AuthUseCase) Logout(ctx *gin.Context) error {
	userIdStr := ctx.GetString("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	return a.authRepo.SaveAccessLog(ctx, &model.AccessLogs{
		UserID:        userId,
		Action:        "logout",
		IPAddress:     ctx.ClientIP(),
		UserAgent:     ctx.GetHeader("User-Agent"),
		Details:       nil,
		CreatedAt:     time.Now(),
		LogonDatetime: time.Now(),
		LogonResult:   "success",
	})
}

// Admin login for testing purposes
func (a *AuthUseCase) loginAsAdmin(ctx *gin.Context, username string) (*model.LoginResponse, error) {
	user, err := a.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	accesstoken, err := a.GenerateToken(24*time.Hour, &auth.Metadata{
		UserId:    user.ID,
		Name:      user.Name,
		CenterId:  user.Center.ID,
		SectionId: user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.GenerateToken(3*24*time.Hour, &auth.Metadata{
		UserId:    user.ID,
		Name:      user.Name,
		CenterId:  user.Center.ID,
		SectionId: user.SectionID,
	})
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accesstoken,
		RefreshToken: refreshToken,
	}, nil
}
