package middleware

import (
	"context"
	"github.com/tumbleweedd/svc/auth_service/internal/domain/valueobjects"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"net/http"
	"strings"

	authtoken "github.com/tumbleweedd/svc/auth_service/pkg/auth/token"
)

type ctxKeyUserID int
type ctxKeyUserLogin int
type ctxKeyUserRole int

const (
	CtxKeyUserID    ctxKeyUserID    = 0
	CtxKeyUserLogin ctxKeyUserLogin = 0
	CtxKeyUserRole  ctxKeyUserRole  = 0
)

func MWAccessTokenValidator(log logger.Logger, jwtManager authtoken.JWTManagerI) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get access token
			accToken := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))

			// Validate access token and get auth info
			authInfo, err := jwtManager.VerifyAccessToken(accToken)
			if err != nil {
				log.ErrorContext(ctx, "Access token isn't valid", log.String("error", err.Error()))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Set auth context to context
			newCtx := context.WithValue(ctx, CtxKeyUserID, authInfo.UserID)
			newCtx = context.WithValue(newCtx, CtxKeyUserLogin, authInfo.Login)
			newCtx = context.WithValue(newCtx, CtxKeyUserRole, authInfo.Role)

			// Set auth info to logger
			userLogger := log.With(
				log.String("user_id", authInfo.UserID),
				log.String("user_login", authInfo.Login),
				log.String("user_role", authInfo.Role),
			)

			userLogger.InfoContext(newCtx, "User authenticated")

			next.ServeHTTP(w, r.WithContext(newCtx))
		}

		return http.HandlerFunc(fn)
	}
}

func MWAuthorization(log logger.Logger, requiredRoles ...valueobjects.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userRole, ok := ctx.Value(CtxKeyUserRole).(string)
			if !ok || userRole == "" {
				log.ErrorContext(ctx, "User role not found in context")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			for _, role := range requiredRoles {
				if userRole == string(role) {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.WarnContext(ctx, "User does not have the required permissions", log.String("user_role", userRole))
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

		return http.HandlerFunc(fn)
	}
}
