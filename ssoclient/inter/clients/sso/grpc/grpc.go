package ssoclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	v1 "github.com/formangloria83/protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api v1.AuthClient
}

func New(
	ctx context.Context,
	addr string,
	log *slog.Logger,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"
	// TODO
	// Add secure connection

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		), // Closing parenthesis moved here
	)

	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", op, err)
	}

	return &Client{
		api: v1.NewAuthClient(cc),
	}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) IsAdmin(ctx context.Context, usedID int64) (bool, error) {
	const op = "grpc.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &v1.IsAdminRequest{UserId: usedID})
	if err != nil {
		return false, fmt.Errorf("failed to call IsAdmin: %w", op, err)
	}
	return resp.IsAdmin, nil
}
