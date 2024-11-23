package ssogrpc

import (
	"context"
	"errors"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	ssov1 "github.com/tendze/gRPC_AuthService_Proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "sso.grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(interceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
		log: log,
	}, nil
}

func (c *Client) ValidateToken(ctx context.Context, token string, appID int) (userID int, email string, isValid bool, errw error) {
	const op = "sso.grpc.ValidateToken"

	resp, err := c.api.ValidateToken(ctx, &ssov1.ValidateTokenRequest{
		Token: token,
		AppId: int32(appID),
	})
	if err != nil {
		userID, email, isValid = 0, "", false
		grpcStatus, _ := status.FromError(err)
		switch grpcStatus.Code() {
		case codes.InvalidArgument:
			errw = fmt.Errorf("%s: %w", op, errors.New("invalid authorization token"))
			return
		case codes.NotFound:
			errw = fmt.Errorf("%s: %w", op, errors.New("app not found"))
			return
		case codes.Internal:
			errw = fmt.Errorf("%s: %w", op, errors.New("auth service internal error"))
			return
		}
		errw = fmt.Errorf("%s: %w", op, err)
		return
	}
	userID, _ = strconv.Atoi(resp.UserId)
	email, isValid, errw = resp.Email, true, nil
	return
}

// SSOMiddleware SSO Interceptor
func SSOMiddleware(authClient *Client, appID int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		const op = "SSOMiddleware"
		log := authClient.log.With(
			slog.String("op", op),
		)

		log.Info("extracting authorization token from context")
		token, err := extractTokenFromContext(ctx)
		if err != nil {
			return nil, err
		}

		userID, email, isValid, err := authClient.ValidateToken(ctx, token, appID)
		_ = userID // user_id is unnecessary
		if err != nil {
			return nil, fmt.Errorf("failed to validate token: %w", err)
		}
		if !isValid {
			return nil, errors.New("invalid authorization token")
		}

		ctx = context.WithValue(ctx, "creatorEmail", email)
		return handler(ctx, req)
	}
}

func interceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// extractTokenFromContext извлекает токен из метаданных gRPC-запроса.
func extractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata in context")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.New("missing authorization header")
	}

	// Извлекаем токен из заголовка: "Bearer <token>"
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == "" {
		return "", errors.New("invalid token format")
	}

	return token, nil
}
