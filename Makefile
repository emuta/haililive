pb:
	protoc --go_out=plugins=grpc,paths=source_relative:./ ./protobuf/haililive/haililive.proto

server: export POSTGRES_URL=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable&application_name=haililive
server:
	@echo "Set POSTGRES_URL: "${POSTGRES_URL}
	@go run ./cmd/server

client: export AMQP_URL=amqp://guest:guest@rabbitmq:5672/tornado
client: export GRPC_SERVER_ADDR=localhost:13721
client:
	@go run ./cmd/watcher

clean:
	rm ./protobuf/haililive/haililive.pb.go