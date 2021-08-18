build:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --micro_out=. --go-grpc_opt=paths=source_relative \
    	proto/vessel/vessel.proto

		docker build -t vessel-service .
run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns vessel-service