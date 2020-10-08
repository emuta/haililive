package interceptor

import (
	"context"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LogginInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	// debug request
	logger := log.WithField("method", info.FullMethod)
	logger.Infof("Request: %v", req)

	// debug metadata from incoming context
	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		log.Error("unable to get metadata from incoming context")
	} else {
		logger.Infof("%#v", md)
	}

	resp, err := handler(ctx, req)

	// debug response
	logger.Debugf("Response: %v", resp)
	return resp, err
}

func RecoverInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	var resp interface{}
	var err error
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic error: %v", r)
		}
	}()
	resp, err = handler(ctx, req)
	return resp, err
}
