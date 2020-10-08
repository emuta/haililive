package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"

	"github.com/emuta/haililive/interceptor"
	pb "github.com/emuta/haililive/protobuf/haililive"
	"github.com/emuta/haililive/service"
)

const (
	// Server default port
	DEFAULT_SERVER_PORT = 13721

	// Logger default config
	DEFAULT_LOGGER_DATETIME_FORMAT = "2006-01-02 15:04:05.000000"
	DEFAULT_LOGGER_LEVEL           = "debug"
)

var (
	port     int
	logLevel string
	db       *gorm.DB
)

func init() {
	flag.IntVar(&port, "port", DEFAULT_SERVER_PORT, "server listening port")
	flag.StringVar(&logLevel, "loglevel", DEFAULT_LOGGER_LEVEL, "Log level")
	flag.Parse()

	setupLogger()
	replaceGRPClogger(log.StandardLogger())

	// setup db
	setupPostgres()
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panic(err)
	}
	defer lis.Close()
	log.WithField("addr", lis.Addr()).Info("Server running")
	runServer(lis)
}

func runServer(l net.Listener) {
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			interceptor.RecoverInterceptor,
			interceptor.LogginInterceptor,
		),
	}

	s := grpc.NewServer(opts...)
	// add health status service
	healthgrpc.RegisterHealthServer(s, health.NewServer())
	pb.RegisterHaiLiLiveServiceServer(s, service.NewHaiLiLiveServiceServerImpl(db))

	if err := s.Serve(l); err != nil {
		log.Panic(err)
	}
}

func setupLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: DEFAULT_LOGGER_DATETIME_FORMAT,
		ForceColors:     true,
	})

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	log.WithFields(log.Fields{
		"timestamp_format": DEFAULT_LOGGER_DATETIME_FORMAT,
		"logger_level":     log.GetLevel(),
	}).Infof("Setup logger")
}

func replaceGRPClogger(logger *log.Logger) {
	entry := log.NewEntry(logger)
	grpc_logrus.ReplaceGrpcLogger(entry)
}

func setupPostgres() {
	url, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		log.Fatal("POSTGRES_URL not found in environment")
	}

	var err error
	db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Postgres server connected")
}
