FROM emuta/golang:grpc as build

WORKDIR /go/src/

RUN go get -d github.com/emuta/haililive/cmd/server \
    && cd github.com/emuta/haililive \
    && go build -ldflags="-s -w" -o /root/server ./cmd/server/main.go && upx /root/server

FROM alpine
COPY --from=build /root/server /srv/

EXPOSE 13721

CMD ["/srv/server"]