//go:generate protoc --go_out=. --go-grpc_out=. ../../../api/SystemStatsService.proto --proto_path=../../../api

package internalgrpc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/pipeline"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Server struct {
	UnimplementedSystemStatsStreamServiceServer
	host    string
	port    string
	grpcSrv *grpc.Server
}

func NewServer(host string, port string) *Server {
	server := &Server{
		grpcSrv: grpc.NewServer(),
		host:    host,
		port:    port,
	}
	RegisterSystemStatsStreamServiceServer(server.grpcSrv, server)

	return server
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("server started")
	if err := s.grpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcSrv.GracefulStop()
}

func (s Server) FetchResponse(message *RequestMessage, server SystemStatsStreamService_FetchResponseServer) error {
	log.Printf("fetch response for N = %d and M = %d", message.N, message.M)

	responseTicker := time.NewTicker(time.Duration(message.N) * time.Second)
	done := make(chan bool)
	in := make(pipeline.Bi)
	stages := pipeline.GetStages()

	go func() {
		for {
			select {
			case <-responseTicker.C:
				stat := app.SystemStats{
					ID:          uuid.New(),
					CollectedAt: time.Now(),
				}
				in <- stat
			case <-done:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			default:
				for stat := range pipeline.ExecutePipeline(in, nil, stages...) {
					log.Printf("Stat %s", stat)

					resp := ResponseMessage{
						Title: fmt.Sprintf("Request for N = %d and M = %d: %s", message.N, message.M, stat),
					}

					if err := server.Send(&resp); err != nil {
						log.Printf("send error %s", err)
						done <- true
						return
					}

					log.Printf("finishing request number")
				}
			}
		}
	}()

	<-done

	return nil
}
