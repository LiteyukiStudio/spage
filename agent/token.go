package agent

import (
	"context"
	"github.com/LiteyukiStudio/spage/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// UnaryInterceptor 鉴权
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 从上下文提取Metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	// 检查authorization头
	auths, ok := md["authorization"]
	if !ok || len(auths) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	// 提取Token
	token := strings.TrimPrefix(auths[0], "Bearer ")
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
	}

	// 验证Token
	if err := ValidateToken(token); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Token有效，继续处理请求
	return handler(ctx, req)
}

// ValidateToken 验证JWT令牌
func ValidateToken(token string) error {
	if token != config.AgentConfig.Server.Token {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return nil
}
