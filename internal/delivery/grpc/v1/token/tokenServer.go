package token

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tokenProto "github.com/tumbleweedd/svc/auth_service/gen/token"
)

func (s *Server) Login(ctx context.Context, req *tokenProto.LoginRequest) (*tokenProto.LoginResponse, error) {
	input := NewLoginRequest(req)

	output, err := s.tokenUseCase.GenerateTokens(ctx, &input)
	if err != nil {
		return &tokenProto.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	return NewGenerateTokenRequest(output), nil
}

func (s *Server) Refresh(ctx context.Context, req *tokenProto.RefreshRequest) (*tokenProto.RefreshResponse, error) {
	return &tokenProto.RefreshResponse{}, nil
}

func (s *Server) Logout(ctx context.Context, req *tokenProto.LogoutRequest) (*tokenProto.LogoutResponse, error) {
	return &tokenProto.LogoutResponse{}, nil
}

func (s *Server) Validate(ctx context.Context, req *tokenProto.ValidateRequest) (*tokenProto.ValidateResponse, error) {
	return &tokenProto.ValidateResponse{}, nil
}
