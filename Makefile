.PHONEY: gen_proto

gen_proto:
	export PATH="$$PATH:$$(go env GOPATH)/bin" && \
  	echo "Generating protobuf code for parser..."  && \
	protoc --go_out=$$PWD/services/parser/internal/api/proto --go-grpc_out=$$PWD/services/parser/internal/api/proto protos/*.proto && \
	echo "Generating protobuf code for rest..."  && \
    protoc --go_out=$$PWD/services/core/internal/grpc/proto --go-grpc_out=$$PWD/services/core/internal/grpc/proto protos/*.proto && \
	echo "Protobuf code generation completed."

up:
	docker-compose up kafka -d && \
	docker compose exec -w /opt/kafka/bin kafka sh ./kafka-topics.sh --bootstrap-server localhost:9092 \
	--create --if-not-exists \
    --topic raw-messages \
    --partitions 3 \
    --replication-factor 1 \
    --config retention.ms=604800000 \
    --config retention.bytes=1073741824 \
    --config cleanup.policy=delete \
    --config max.message.bytes=10485760
upBuild:
	docker-compose up -d --build
down:
	docker-compose down
