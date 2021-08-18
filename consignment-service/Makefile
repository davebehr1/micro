
build:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --micro_out=. --go-grpc_opt=paths=source_relative \
    	proto/consignment/consignment.proto

	docker build -t consignment-service .

run: 
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns consignment-service