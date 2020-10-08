package main

import (
	"context"
	"flag"
	"os"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/emuta/haililive/protobuf/haililive"
	"github.com/emuta/haililive/watcher"
	"github.com/emuta/haililive/watcher/handler"
)

const (
	DEFAULT_LOGGER_LEVEL           = "debug"
	DEFAULT_LOGGER_DATETIME_FORMAT = "2006-01-02 15:04:05.000000"

	DEFAULT_AMQP_EXCHANGE_NAME = "stock"
)

var (
	logLevel           string
	GRPC_SERVER_ADDR   string
	AMQP_URL           string
	AMQP_EXCHANGE_NAME string
	AMQP_QUEUE_NAME    string = os.Getenv("AMQP_QUEUE_NAME")

	cli pb.HaiLiLiveServiceClient
)

func init() {
	flag.StringVar(&logLevel, "loglevel", DEFAULT_LOGGER_LEVEL, "Log level")
	flag.Parse()

	setupLogger()
	replaceGRPClogger(log.StandardLogger())

	initGRPCclient()
	initAMQPcfg()
}

func main() {
	ctx := context.Background()
	h := handler.NewHandler(cli)
	w := watcher.NewWatcher(AMQP_URL, AMQP_EXCHANGE_NAME, AMQP_QUEUE_NAME, h)
	log.WithFields(log.Fields{
		"exchange": AMQP_EXCHANGE_NAME,
		"queue":    AMQP_QUEUE_NAME,
		"grpc":     GRPC_SERVER_ADDR,
	}).Info("config")
	w.Run(ctx)
}

func initAMQPcfg() {
	if v, ok := os.LookupEnv("AMQP_URL"); ok {
		AMQP_URL = v
	} else {
		log.Fatal("AMQP_URL not found in environment")
	}

	if v, ok := os.LookupEnv("AMQP_EXCHANGE_NAME"); ok {
		AMQP_EXCHANGE_NAME = v
	} else {
		AMQP_EXCHANGE_NAME = DEFAULT_AMQP_EXCHANGE_NAME
	}
}

func initGRPCclient() {

	if v, ok := os.LookupEnv("GRPC_SERVER_ADDR"); ok {
		GRPC_SERVER_ADDR = v
	} else {
		log.Fatal("GRPC_SERVER_ADDR not found in environment")
	}

	opt := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	cc, err := grpc.Dial(GRPC_SERVER_ADDR, opt...)
	if err != nil {
		log.Fatal(err)
	}

	cli = pb.NewHaiLiLiveServiceClient(cc)
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
