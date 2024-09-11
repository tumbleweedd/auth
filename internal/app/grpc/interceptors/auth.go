package interceptors

import (
	"context"
	"github.com/tumbleweedd/svc/auth_service/internal/domain/valueobjects"
	"strings"

	tokenProto "github.com/tumbleweedd/svc/auth_service/gen/token"
	userProto "github.com/tumbleweedd/svc/auth_service/gen/user"
	authtoken "github.com/tumbleweedd/svc/auth_service/pkg/auth/token"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager authtoken.JWTManagerI
}

func NewAuthInterceptor(jwtManager authtoken.JWTManagerI) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
	}
}

func (a *AuthInterceptor) UnaryAuthentication() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == userProto.UserService_Create_FullMethodName ||
			info.FullMethod == tokenProto.TokenService_Login_FullMethodName ||
			info.FullMethod == tokenProto.TokenService_Refresh_FullMethodName {

			return handler(ctx, req)
		}

		// Проверяем аутентификацию
		authInfo, err := a.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		// Добавляем claims в контекст, чтобы его могли использовать дальнейшие обработчики
		newCtx := context.WithValue(ctx, "authInfo", authInfo)
		return handler(newCtx, req)
	}
}

func (a *AuthInterceptor) authenticate(ctx context.Context) (*authtoken.AuthInfo, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	var tokenString string
	if authHeaders := md["Authorization"]; len(authHeaders) > 0 {
		tokenString = authHeaders[0]
	}

	if tokenString == "" {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// Извлекаем токен из строки авторизации
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token format")
	}

	tokenString = parts[1]

	authInfo, err := a.jwtManager.VerifyAccessToken(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "invalid token")
	}

	return authInfo, nil
}

var noAuthAccessMap = map[string]map[valueobjects.Role]bool{
	userProto.UserService_Create_FullMethodName: {
		valueobjects.UserRole:  true,
		valueobjects.AdminRole: true,
	},
	tokenProto.TokenService_Login_FullMethodName: {
		valueobjects.AdminRole: true,
		valueobjects.UserRole:  true,
	},
	tokenProto.TokenService_Refresh_FullMethodName: {
		valueobjects.AdminRole: true,
		valueobjects.UserRole:  true,
	},
}

func (a *AuthInterceptor) UnaryAuthorization() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == userProto.UserService_Create_FullMethodName ||
			info.FullMethod == tokenProto.TokenService_Login_FullMethodName ||
			info.FullMethod == tokenProto.TokenService_Refresh_FullMethodName {

			return handler(ctx, req)
		}

		// Проверяем авторизацию
		authInfo := ctx.Value("authInfo").(*authtoken.AuthInfo)
		if err := a.authorize(authInfo, info.FullMethod); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) authorize(authInfo *authtoken.AuthInfo, method string) error {
	value, ok := noAuthAccessMap[method][valueobjects.Role(authInfo.Role)]
	if ok && value {
		return nil
	}

	return status.Errorf(codes.PermissionDenied, "permission denied")
}
