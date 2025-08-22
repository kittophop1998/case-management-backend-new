package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Metadata struct {
	Name       string    `json:"name"`
	UserId     uuid.UUID `json:"userId"`
	CenterId   uuid.UUID `json:"centerId"`
	CenterName string    `json:"centerName"`
	SectionId  uuid.UUID `json:"sectionId"`
}

type JwtClaims struct {
	jwt.StandardClaims
	Metadata
}
