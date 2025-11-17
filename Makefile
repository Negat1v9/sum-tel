.PHONEY: gen_proto

gen_proto:
	export PATH="$$PATH:$$(go env GOPATH)/bin" && \
  	echo "Generating protobuf code for parser..."  && \
	protoc --go_out=$$PWD/services/parser/internal/grpc/proto --go-grpc_out=$$PWD/services/parser/internal/grpc/proto protos/*.proto && \
	echo "Generating protobuf code for rest..."  && \
    protoc --go_out=$$PWD/services/core/internal/grpc/proto --go-grpc_out=$$PWD/services/core/internal/grpc/proto protos/*.proto && \
	echo "Protobuf code generation completed."

up:
	docker-compose up -d
upBuild:
	docker-compose up -d --build
down:
	docker-compose down