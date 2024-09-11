package token

import (
	tokenProto "github.com/tumbleweedd/svc/auth_service/gen/token"
	"github.com/tumbleweedd/svc/auth_service/internal/domain/usecase/token"
	"google.golang.org/grpc"
)

type Server struct {
	tokenUseCase *token.UseCase

	tokenProto.UnimplementedTokenServiceServer
}

func RegisterServer(
	gRPCServer *grpc.Server,
	tokenUC *token.UseCase,
) {
	tokenProto.RegisterTokenServiceServer(gRPCServer, &Server{tokenUseCase: tokenUC})
}
