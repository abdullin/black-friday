.PHONY: clean schema


schema:
	protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:.  \
    	protos/ch1.proto


	protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:.  \
    	seq/test.proto
