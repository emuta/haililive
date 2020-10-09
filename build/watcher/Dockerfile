FROM emuta/golang:grpc as build

WORKDIR /go/src/

RUN go get -d github.com/emuta/haililive/cmd/watcher \
    && cd github.com/emuta/haililive \
    && go build -ldflags="-s -w" -o /root/watcher ./cmd/watcher/main.go && upx /root/watcher

FROM alpine

ENV TZ Asia/Shanghai

WORKDIR /srv

COPY --from=build /root/watcher ./

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk update \
    && apk add tzdata --no-cache \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && rm -rf /var/cache/apk


CMD ["./watcher"]