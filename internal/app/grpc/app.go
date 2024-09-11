package grpc

import (
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/tumbleweedd/svc/auth_service/internal/app/grpc/interceptors"
	tokenGRPC "github.com/tumbleweedd/svc/auth_service/internal/delivery/grpc/v1/token"
	tokenUseCase "github.com/tumbleweedd/svc/auth_service/internal/domain/usecase/token"
	userUseCase "github.com/tumbleweedd/svc/auth_service/internal/domain/usecase/user"
	"github.com/tumbleweedd/svc/auth_service/pkg/auth/token"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        logger.Logger
	jwtManager token.JWTManagerI
	gRPCServer *grpc.Server
	port       int
}

func New(
	log logger.Logger,
	jwtManager token.JWTManagerI,
	tokenUC *tokenUseCase.UseCase,
	userUC *userUseCase.UseCase,
	port int,
) *App {
	recoveryOptions := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("panic triggered: %v", p)

			return status.Errorf(codes.Internal, "panic triggered: %v", p)
		}),
	}

	authInterceptor := interceptors.NewAuthInterceptor(jwtManager)
	loggerInterceptor := interceptors.NewLoggerInterceptor(log)

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOptions...),

		loggerInterceptor.UnaryLogger(),
		authInterceptor.UnaryAuthentication(),
		authInterceptor.UnaryAuthorization(),
	))

	tokenGRPC.RegisterServer(gRPCServer, tokenUC)

	return &App{
		log:        log,
		jwtManager: jwtManager,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(fmt.Sprintf("failed to run gRPC server: %v", err))
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info(fmt.Sprintf("op %s: starting gRPC server on port %d", op, a.port))

	return a.gRPCServer.Serve(listener)
}
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.Info(fmt.Sprintf("op %s: stopping gRPC server on port %d", op, a.port))

	a.gRPCServer.GracefulStop()
}
