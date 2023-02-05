
bin/bf:
	@go build -buildmode=pie -o bin/bf main.go

.PHONY: schema
schema:
	protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:.  \
    	inventory/api/api.proto

.PHONY: clean
clean:
	@rm -rf bin/bf

.PHONY: test
test: bin/bf
	@bin/bf test

.PHONY: perf
perf: bin/bf
	@bin/bf perf

.PHONY: stress
stress: bin/bf
	@bin/bf stress
