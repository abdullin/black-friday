.PHONY: clean schema specs


schema:
	protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:.  \
    	api/api.proto


specs:
	@go run  -buildmode=pie *.go
