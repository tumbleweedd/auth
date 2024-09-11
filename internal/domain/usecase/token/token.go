package token

import (
	"context"
	"github.com/tumbleweedd/svc/auth_service/pkg/auth/token"

	tokenEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/token"
)

type tokenService interface {
	GenerateTokens(
		ctx context.Context,
		login, password string,
	) (accessToken *token.TokenInfo, refreshToken *token.TokenInfo, err error)
	RefreshToken(
		ctx context.Context,
		refreshToken string,
	) (refreshedToken *token.TokenInfo, err error)
}

type UseCase struct {
	tokenService tokenService
}

func NewUseCase(tokenService tokenService) *UseCase {
	return &UseCase{tokenService: tokenService}
}

func (u *UseCase) GenerateTokens(
	ctx context.Context,
	login, password string,
) (accessToken, refreshToken *tokenEntity.Token, err error) {
	accessT, refreshT, err := u.tokenService.GenerateTokens(ctx, login, password)
	if err != nil {
		return nil, nil, err
	}

	return &tokenEntity.Token{
			Token:     accessT.Token,
			IssuedAt:  accessT.IssuedAt,
			ExpiresAt: accessT.ExpiresAt,
		}, &tokenEntity.Token{
			Token:     refreshT.Token,
			IssuedAt:  refreshT.IssuedAt,
			ExpiresAt: refreshT.ExpiresAt,
		}, nil
}
