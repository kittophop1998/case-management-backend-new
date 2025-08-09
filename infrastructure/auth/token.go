package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Metadata struct {
	Name      string    `json:"name"`
	UserId    uuid.UUID `json:"userId"`
	CenterId  uuid.UUID `json:"centerId"`
	SectionId uuid.UUID `json:"sectionId"`
	QueueId   uuid.UUID `json:"queueId"`
}

type JwtClaims struct {
	jwt.StandardClaims
	Metadata
}
