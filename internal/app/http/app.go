package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tumbleweedd/svc/auth_service/internal/app/http/middleware"
	"github.com/tumbleweedd/svc/auth_service/internal/config"
	createUserHTTPHandler "github.com/tumbleweedd/svc/auth_service/internal/delivery/http/v1/user/create"
	deleteUserHTTPHandler "github.com/tumbleweedd/svc/auth_service/internal/delivery/http/v1/user/delete"
	getUserHTTPHandler "github.com/tumbleweedd/svc/auth_service/internal/delivery/http/v1/user/get"
	"github.com/tumbleweedd/svc/auth_service/internal/domain/valueobjects"
	"github.com/tumbleweedd/svc/auth_service/pkg/auth/token"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"net/http"
)

type App struct {
	log         logger.Logger
	httpServer  *http.Server
	jwtManager  token.JWTManagerI
	userCreator createUserHTTPHandler.UserCreator
	userRemover deleteUserHTTPHandler.UserRemover
	userGetter  getUserHTTPHandler.UserGetter
}

func NewApp(
	log logger.Logger,
	cfg *config.HTTPConfig,
	jwtManager token.JWTManagerI,
	userCreator createUserHTTPHandler.UserCreator,
	userRemover deleteUserHTTPHandler.UserRemover,
	userGetter getUserHTTPHandler.UserGetter,
) *App {
	mux := chi.NewRouter()

	// Регистрация pprof-обработчиков
	//mux.HandleFunc("/debug/pprof/", pprof.Index)
	//mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	createH := createUserHTTPHandler.NewHandler(log, userCreator)
	getH := getUserHTTPHandler.NewHandler(log, userGetter)
	removeH := deleteUserHTTPHandler.NewHandler(log, userRemover)

	mux.Route("/v1", func(r chi.Router) {
		// No auth
		r.Group(func(r chi.Router) {
			// User
			r.Post("/users", createH.Create)
		})

		// Auth
		r.Group(func(r chi.Router) {
			r.Use(middleware.MWAccessTokenValidator(log, jwtManager))
			r.Use(middleware.MWAuthorization(log, valueobjects.AdminRole, valueobjects.UserRole))

			// User
			r.Get("/users", getH.Get)
			r.Delete("/users", removeH.Delete)
		})
	})

	httpServer := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}

	return &App{
		log:        log,
		httpServer: httpServer,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(fmt.Sprintf("failed to run http server: %v", err))
	}
}

func (a *App) run() error {
	a.log.Info("starting http server", a.log.String("port", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("failed to run http server", a.log.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("shutting down http server", a.log.String("port", a.httpServer.Addr))

	return a.httpServer.Shutdown(ctx)
}
