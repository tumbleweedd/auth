package app

import (
	"context"
	"fmt"
	eventsender "github.com/tumbleweedd/svc/auth_service/internal/domain/eventSender/service"
	outboxservice "github.com/tumbleweedd/svc/auth_service/internal/domain/outboxer/service"
	outboxrepo "github.com/tumbleweedd/svc/auth_service/internal/infrastructure/repository/postgres/outbox"
	"github.com/tumbleweedd/svc/auth_service/pkg/auth/token"
	"github.com/tumbleweedd/svc/auth_service/pkg/client/kafka"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	grpcapp "github.com/tumbleweedd/svc/auth_service/internal/app/grpc"
	httpapp "github.com/tumbleweedd/svc/auth_service/internal/app/http"

	"github.com/tumbleweedd/svc/auth_service/internal/config"
	userusecase "github.com/tumbleweedd/svc/auth_service/internal/domain/usecase/user"
	userservice "github.com/tumbleweedd/svc/auth_service/internal/domain/user/service"
	userinmemoryrepo "github.com/tumbleweedd/svc/auth_service/internal/infrastructure/repository/cacheImpl/user"
	userrepo "github.com/tumbleweedd/svc/auth_service/internal/infrastructure/repository/postgres/user"
	usersercretrepo "github.com/tumbleweedd/svc/auth_service/internal/infrastructure/repository/postgres/userSercret"
	"github.com/tumbleweedd/svc/auth_service/pkg/client/postgres"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
)

type App struct {
	gRPCServer *grpcapp.App
	hTTPServer *httpapp.App
	postgresDB *postgres.PgDB
}

func New(
	ctx context.Context,
	cfg *config.Config,
	log logger.Logger,
	port int,
) (
	a *App,
) {
	pgDB, err := postgres.NewClient(ctx, log, &cfg.PgConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to postgres: %v", err))
	}

	userRepository := userrepo.NewUserRepository(pgDB.GetDB(), log, trmsqlx.DefaultCtxGetter)
	userSecretRepository := usersercretrepo.NewUserSecretRepository(pgDB.GetDB(), log, trmsqlx.DefaultCtxGetter)
	userInMemoryRepository := userinmemoryrepo.NewInMemoryUserRepository(ctx, cfg.CacheConfig.DefaultTTL)
	outboxRepository := outboxrepo.NewOutboxRepository(pgDB.GetDB(), log, trmsqlx.DefaultCtxGetter)

	trManager := manager.Must(trmsqlx.NewDefaultFactory(pgDB.GetDB()))

	userService := userservice.NewService(
		userRepository,
		userSecretRepository,
		userInMemoryRepository,
		trManager,
		log,
	)

	outboxService := outboxservice.NewService(outboxRepository)

	userUC := userusecase.NewUseCase(userService, outboxService)

	jwtManager := token.NewJWTManager(
		cfg.JWTManagerConfig.AccessSecret,
		cfg.JWTManagerConfig.RefreshSecret,
		cfg.JWTManagerConfig.AccessTimeout,
		cfg.JWTManagerConfig.RefreshTimeout,
	)

	gRPCApp := grpcapp.New(log, jwtManager, nil, userUC, port)
	httpApp := httpapp.NewApp(
		log,
		&cfg.HTTPConfig,
		jwtManager,
		userUC,
		userUC,
		userUC,
	)

	kafkaProducer := kafka.NewProducer(cfg.KafkaConfig.Port, log)

	sender := eventsender.New(outboxRepository, outboxRepository, &cfg.KafkaConfig, kafkaProducer, log)
	sender.StartProcessingEvents(ctx, cfg.EventSenderConfig.HandlePeriodMin*time.Minute)

	return &App{
		gRPCServer: gRPCApp,
		hTTPServer: httpApp,
		postgresDB: pgDB,
	}
}

func (a *App) MustRun() {
	go a.gRPCServer.MustRun()
	go a.hTTPServer.MustRun()
}

func (a *App) Stop(ctx context.Context) {
	a.gRPCServer.Stop()

	err := a.hTTPServer.Shutdown(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to shutdown http server: %v", err))
	}

	err = a.postgresDB.Close()
	if err != nil {
		panic(fmt.Sprintf("failed to close postgres: %v", err))
	}
}
