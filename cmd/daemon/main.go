package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/vlyagusha/system_stats_daemon/internal/config"
	internalgrpc "github.com/vlyagusha/system_stats_daemon/internal/server/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "port", "50005", "port")
}

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	grpcServer := internalgrpc.NewServer("", port, *cfg)

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("failed to start grpc server: %s", err)
		}
	}()

	<-ctx.Done()
	log.Printf("graceful shutting down")
	grpcServer.Stop()
}
