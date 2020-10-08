FROM emuta/golang:grpc as build

WORKDIR /go/src/

RUN go get -d github.com/emuta/haililive/cmd/server \
    && cd github.com/emuta/haililive \
    && go build -ldflags="-s -w" -o /root/server ./cmd/server/main.go && upx /root/server

FROM scratch
COPY --from=build /root/server /

EXPOSE 13721

CMD ["/server"]