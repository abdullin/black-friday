.PHONY: clean schema


schema:
	protoc --go_out=paths=source_relative:./sdk-go \
		--go-grpc_out=paths=source_relative:./sdk-go  \
    	protos/ch1.proto
	python -m grpc_tools.protoc -Iprotos --python_out=src-qa/qa/protos --grpc_python_out=src-qa/qa/protos protos/ch1.proto
