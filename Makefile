.PHONY: clean schema


build:
	@go build -buildmode=pie -o bin/bf main.go



schema:
	protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:.  \
    	inventory/api/api.proto

clean:
	@rm -rf bin/bf


test: build
	@bin/bf test
