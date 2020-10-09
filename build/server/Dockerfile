FROM emuta/golang:grpc as build

WORKDIR /go/src/

RUN go get -d github.com/emuta/haililive/cmd/server \
    && cd github.com/emuta/haililive \
    && go build -ldflags="-s -w" -o /root/server ./cmd/server/main.go && upx /root/server


FROM alpine

ENV TZ Asia/Shanghai

WORKDIR /srv

COPY --from=build /root/server ./

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk update \
    && apk add tzdata --no-cache \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && rm -rf /var/cache/apk

EXPOSE 13721

CMD ["./server"]