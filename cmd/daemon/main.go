package main

import (
	"context"
	"flag"
	internalgrpc "github.com/vlyagusha/system_stats_daemon/internal/server/grpc"
	"log"
	"os/signal"
	"syscall"
)

var port string

func init() {
	flag.StringVar(&port, "port", "50005", "port")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	grpcServer := internalgrpc.NewServer("", port)

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("failed to start grpc server: %s", err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Printf("graceful shutting down")
		grpcServer.Stop()
	}()

	<-ctx.Done()
}
