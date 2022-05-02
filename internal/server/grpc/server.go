//go:generate protoc --go_out=. --go-grpc_out=. ../../../api/SystemStatsService.proto --proto_path=../../../api

package internalgrpc

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/vlyagusha/system_stats_daemon/internal/app"
	"github.com/vlyagusha/system_stats_daemon/internal/config"
	"github.com/vlyagusha/system_stats_daemon/internal/pipeline"
	memorystorage "github.com/vlyagusha/system_stats_daemon/internal/storage/memory"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedSystemStatsStreamServiceServer
	host    string
	port    string
	grpcSrv *grpc.Server
	config  config.Config
}

func NewServer(host string, port string, config config.Config) *Server {
	server := &Server{
		grpcSrv: grpc.NewServer(),
		host:    host,
		port:    port,
		config:  config,
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

	in := make(pipeline.Bi)
	done := make(chan bool)
	collectDone := make(chan bool)

	collectTicker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-collectTicker.C:
				stat := app.SystemStats{
					ID:          uuid.New(),
					CollectedAt: time.Now(),
				}
				in <- stat
				log.Printf("pushed stat: %s", stat)
			case <-collectDone:
				return
			}
		}
	}()

	stages := pipeline.GetStages(s.config.Stats)
	storage := memorystorage.New()
	go func(storage *memorystorage.Storage) {
		for stat := range pipeline.ExecutePipeline(in, nil, stages...) {
			err := storage.Create(stat)
			log.Printf("stored stat: %s", stat)
			if err != nil {
				log.Printf("error while collect stats: %s", err)
				return
			}
		}
	}(storage)

	responseTicker := time.NewTicker(time.Duration(message.N) * time.Second)
	go func(storage *memorystorage.Storage) {
		if message.M > message.N {
			log.Printf("sleep %d seconds", message.M-message.N)
			time.Sleep(time.Duration(message.M-message.N) * time.Second)
		}

		for {
			select {
			case <-responseTicker.C:
				stat, err := storage.FindAvg(time.Duration(message.M))
				if err != nil {
					log.Printf("error while getting avg stats: %s", err)
					collectDone <- true
					done <- true
					return
				}

				log.Printf("got avg stat %s", stat)
				resp := ResponseMessage{
					Title:       fmt.Sprintf("Request for N = %d and M = %d: %s", message.N, message.M, stat),
					CollectedAt: time.Now().String(),
					Load: &LoadMessage{
						Load1:  float32(stat.Load1),
						Load5:  float32(stat.Load5),
						Load15: float32(stat.Load15),
					},
					Cpu: &CPUMessage{
						User:   float32(stat.User),
						System: float32(stat.System),
						Idle:   float32(stat.Idle),
					},
					Disk: &DiskMessage{
						Kbt: float32(stat.KBt),
						Tps: float32(stat.TPS),
						Mbs: float32(stat.MBs),
					},
				}

				if err := server.Send(&resp); err != nil {
					log.Printf("send error %s", err)
					collectDone <- true
					done <- true
					return
				}
				log.Printf("finishing request number")

			case <-done:
				log.Printf("finished fetch response for N = %d and M = %d", message.N, message.M)
				return
			}
		}
	}(storage)

	<-done
	log.Printf("finished fetch response for N = %d and M = %d", message.N, message.M)

	return nil
}
