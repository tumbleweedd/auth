package user

import (
	"context"
	userProto "github.com/tumbleweedd/auth_service_proto/gen/go/user"
	"github.com/tumbleweedd/svc/auth_service/internal/domain/usecase/user"
	"google.golang.org/grpc"
)

type userUC interface {
	Create(ctx context.Context, req *user.CreateUserRequest) (user.CreateUserOutput, error)
}

type Server struct {
	userUseCase userUC

	userProto.UnimplementedUserServiceServer
}

func RegisterServer(gRPCServer *grpc.Server, userUseCase userUC) {
	userProto.RegisterUserServiceServer(gRPCServer, &Server{userUseCase: userUseCase})
}
