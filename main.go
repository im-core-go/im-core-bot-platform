package main

import (
	"net"
	"os"

	"github.com/im-core-go/im-core-bot-platform/configs"
	grpcserver "github.com/im-core-go/im-core-bot-platform/internal/grpc"
	"github.com/im-core-go/im-core-bot-platform/internal/svc"
	"github.com/im-core-go/im-core-bot-platform/pkg/logger"
	"github.com/im-core-go/im-core-proto/gen/bot/v1"

	"google.golang.org/grpc"
)

func main() {
	lgr := logger.L()
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/dev.yaml"
	}
	cfg, err := configs.Load(configPath)
	if err != nil {
		lgr.Fatalf("load config error: %v", err)
	}
	svcCtx := svc.NewContext(cfg)
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = ":9090"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		lgr.Fatalf("listen error: %v", err)
	}
	server := grpc.NewServer()
	chatServer, err := grpcserver.NewChatServer(svcCtx)
	if err != nil {
		lgr.Fatalf("grpc server init error: %v", err)
	}
	botv1.RegisterChatServiceServer(server, chatServer)
	lgr.Infof("grpc server start on %s", addr)
	if err := server.Serve(listener); err != nil {
		lgr.Fatalf("grpc server stopped: %v", err)
	}
}
