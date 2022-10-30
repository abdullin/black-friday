.PHONY: clean schema specs


build:
	@go build  -buildmode=pie  -o bin/bf main.go



schema:
	protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:.  \
    	inventory/api/api.proto




specs:
	@go run  -buildmode=pie *.go
