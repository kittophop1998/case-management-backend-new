package handler

import (
	"case-management/infrastructure/auth"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := extractBearerToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized", "details": err.Error()})
			return
		}

		claims, err := parseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized", "details": err.Error()})
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Set("username", claims.Name)
		ctx.Next()
	}
}

func extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

func parseToken(tokenString string) (*auth.JwtClaims, error) {
	secretKey := "casemanagement_secret_key" // Replace with your actual secret key

	token, err := jwt.ParseWithClaims(tokenString, &auth.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("couldn't parse claims or token is invalid")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
