package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

const (
	SecretKey = "daDq2Ra@FsA15D"
)

// token cache
var tokenMap = map[string]string{"123456": "secret_key"}

func Auth(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errors.New("无Token认证信息")
	}

	// authentication
	val := md.Get("app_id")
	if len(val) <= 0 {
		return ctx, errors.New("app_id is empty")
	}
	appID := val[0]

	val = md.Get("secret_key")
	if len(val) <= 0 {
		return ctx, errors.New("secret_key is empty")
	}
	appKey := val[0]

	if _, ok := tokenMap[appID]; !ok || appKey != tokenMap[appID] {
		return ctx, errors.New("unauthorized")
	}

	return ctx, nil
}
