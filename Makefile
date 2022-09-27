.PHONY: clean schema


schema:
	protoc --go_out=paths=source_relative:./sdk-go \
		--go-grpc_out=paths=source_relative:./sdk-go  \
		--mypy_out=paths=source_relative:./src-qa/qa \
		--mypy_grpc_out=paths=source_relative:./src-qa/qa \
    	protos/ch1.proto
	python -m grpc_tools.protoc -Iprotos --python_out=src-qa/qa/protos --grpc_python_out=src-qa/qa/protos protos/ch1.proto
	find src-qa/qa/protos/ -type f -name "*.py" -print0 | xargs -0 sed -i.bak 's,import ch1_pb2,from . import ch1_pb2,g'
	find src-qa/qa/protos/ -type f -name "*grpc.pyi" -print0 | xargs -0 sed -i.bak 's,import protos.ch1_pb2,from qa import protos,g'
	rm src-qa/qa/protos/*.bak
