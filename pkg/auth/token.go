package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Metadata struct {
	Username string    `json:"username"`
	UserId   uuid.UUID `json:"userId"`
	CenterId uuid.UUID `json:"centerId"`
	TeamId   uuid.UUID `json:"teamId"`
	QueueId  uuid.UUID `json:"queueId"`
}

type JwtClaims struct {
	jwt.StandardClaims
	Metadata
}
