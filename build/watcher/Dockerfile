FROM emuta/golang:grpc as build

WORKDIR /go/src/

RUN go get -d github.com/emuta/haililive/cmd/watcher \
    && cd github.com/emuta/haililive \
    && go build -ldflags="-s -w" -o /root/watcher ./cmd/watcher/main.go && upx /root/watcher

FROM scratch
COPY --from=build /root/watcher /


CMD ["/watcher"]