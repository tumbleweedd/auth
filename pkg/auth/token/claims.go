package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	AuthInfo
}

type AuthInfo struct {
	UserID string
	Role   string
	Login  string
}

type TokenInfo struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
