package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	grpcapi "github.com/ubombi/timeseries/api/grpc"
	httpapi "github.com/ubombi/timeseries/api/http"
	"github.com/ubombi/timeseries/storage"

	"github.com/ubombi/timeseries/storage/clickhouse"
	"github.com/ubombi/timeseries/storage/dummy"

	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
)

var (
	grpcListenAddr = flag.String("grpcListenAddr", ":8088", "TCP Address to serve as gRPC")
	httpListenAddr = flag.String("httpListenAddr", ":8008", "TCP Address to serve as HTTP")
	storageName    = flag.String("storage", "clickhouse", "storage engine. 'clickhouse' or ''dummy'")
)

func main() {
	flag.Parse()
	// for gracefull shutdown and cancelation
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(3)
	Storage := initStorage(ctx, wg)
	GRPCServer := initGRPC(ctx, wg, Storage)
	HTTPServer := initHTTP(ctx, wg, Storage)

	log.Print("service started")
	WaitForSigterm()
	log.Print("received shutdown signal")
	go cancel()
	go GRPCServer.GracefulStop()
	if graceful, ok := Storage.(storage.Graceful); ok {
		graceful.Shutdown()
	}
	go func() {
		if err := HTTPServer.Shutdown(); err != nil {
			log.Print(err)
		}
	}()
}

func WaitForSigterm() os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return <-ch
}

func initGRPC(ctx context.Context, wg *sync.WaitGroup, storage storage.Interface) (server *grpc.Server) {
	lis, err := net.Listen("tcp", *grpcListenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)

	}

	server = grpc.NewServer()

	grpcapi.RegisterEventServiceServer(
		server,
		&grpcapi.Server{
			Storage: storage,
		},
	)

	go func() {
		err := server.Serve(lis) // Blocking
		if err != nil {
			log.Print(err)
		}
		log.Print("grpc off")
		wg.Done()
	}()

	return server
}

func initHTTP(ctx context.Context, wg *sync.WaitGroup, storage storage.Interface) *fasthttp.Server {
	service := httpapi.Service{Storage: storage}

	s := &fasthttp.Server{
		Handler: service.HandlerFunc,
		Name:    "nginx/1.15.7",
	}

	go func() {
		if err := s.ListenAndServe(*httpListenAddr); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
		log.Print("http off")
		wg.Done()
	}()

	return s
}

func initStorage(ctx context.Context, wg *sync.WaitGroup) storage.Interface {
	switch *storageName {
	case "clickhouse":
		chstorage := clickhouse.NewStorage(ctx)
		go func() {
			defer wg.Done()
			if err := chstorage.Start(); err != nil {
				log.Fatal(err)
			}
			log.Print("storage off")
		}()
		return chstorage
	case "dummy":
		wg.Done()
		return dummy.Stdout
	default:
		log.Fatal("unsupported storage ", *storageName)
	}
	return nil
}
