package token

import (
	tokenProto "github.com/tumbleweedd/svc/auth_service/gen/token"
	"github.com/tumbleweedd/svc/auth_service/internal/domain/usecase/token"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewLoginRequest(req *tokenProto.LoginRequest) token.GenerateTokensInput {
	return token.GenerateTokensInput{
		Login:    req.Username,
		Password: req.Password,
	}
}

func NewGenerateTokenRequest(req *token.GenerateTokensOutput) *tokenProto.LoginResponse {
	tokenInfo := func(token *token.Token) *tokenProto.TokenInfo {
		return &tokenProto.TokenInfo{
			Token:     token.Token,
			IssuedAt:  timestamppb.New(token.IssuedAt),
			ExpiresAt: timestamppb.New(token.ExpiresAt),
		}
	}

	return &tokenProto.LoginResponse{
		AccessToken:  tokenInfo(&req.AccessToken),
		RefreshToken: tokenInfo(&req.RefreshToken),
	}
}
