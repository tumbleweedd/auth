package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTManagerI interface {
	GenerateAccessToken(userInfo *AuthInfo) (*TokenInfo, error)
	GenerateRefreshToken(userInfo *AuthInfo) (*TokenInfo, error)
	VerifyAccessToken(token string) (*AuthInfo, error)
	VerifyRefreshToken(token string) (*AuthInfo, error)
}

type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte

	accessTimeout  int
	refreshTimeout int
}

func NewJWTManager(
	accessSecret, refreshSecret string,
	accessTimeout, refreshTimeout int,
) *JWTManager {
	return &JWTManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),

		accessTimeout:  accessTimeout,
		refreshTimeout: refreshTimeout,
	}
}

func (jwtManager *JWTManager) GenerateAccessToken(userInfo *AuthInfo) (*TokenInfo, error) {
	return createToken(jwtManager.accessSecret, jwtManager.accessTimeout, userInfo)
}

func (jwtManager *JWTManager) GenerateRefreshToken(userInfo *AuthInfo) (*TokenInfo, error) {
	return createToken(jwtManager.refreshSecret, jwtManager.refreshTimeout, userInfo)
}

func createToken(tokenKey []byte, tokenTimeoutMin int, userInfo *AuthInfo) (*TokenInfo, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Duration(tokenTimeoutMin) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		AuthInfo: *userInfo,
	})

	tokenSigned, err := token.SignedString(tokenKey)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		Token:     tokenSigned,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}

func (jwtManager *JWTManager) VerifyAccessToken(token string) (*AuthInfo, error) {
	return verifyToken(jwtManager.accessSecret, token)
}

func (jwtManager *JWTManager) VerifyRefreshToken(token string) (*AuthInfo, error) {
	return verifyToken(jwtManager.refreshSecret, token)
}

func verifyToken(tokenKey []byte, tokenSigned string) (*AuthInfo, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenSigned, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return tokenKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenSigned)
	}

	return &claims.AuthInfo, nil
}
