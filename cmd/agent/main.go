package main

import (
	"context"
	"fmt"
	"github.com/LiteyukiStudio/spage/agent"
	"github.com/LiteyukiStudio/spage/agent/caddy"
	"github.com/LiteyukiStudio/spage/pkg/config"
	"github.com/sirupsen/logrus"
	"net"
)

func main() {
	// 解析参数
	config.AgentConfig.LoadByCmd()
	err := config.AgentConfig.InitAgentConfig()
	if err != nil {
		logrus.Panicf("Failed to initialize agent configuration: %v", err)
		return
	}
	// 启动caddy
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	caddy.StartDaemon(ctx)
	// 注册 grpc 服务
	grpcServer, err := agent.RegisterGrpcApps()
	if err != nil {
		logrus.Panicf("Failed to register grpc apps: %v", err)
		return
	}

	// 创建 TCP 监听器
	servicePoint := fmt.Sprintf("%s:%s", config.AgentConfig.Service.Host, config.AgentConfig.Service.Port)
	lis, err := net.Listen("tcp", servicePoint)
	if err != nil {
		logrus.Fatalf("Failed to listen on %s: %v", servicePoint, err)
	}

	// 启动 gRPC 服务
	logrus.Printf("gRPC server starting on %s\n", servicePoint)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("gRPC server failed: %v", err)
	}

}
