.PHONY: run fmt check_fmt

run:
	@go run cmd/cli/main.go

test:
	@go test -v ./...

build:
	@cd pkg/lambdaworks/lib/lambdaworks && cargo build --release
	@cp pkg/lambdaworks/lib/lambdaworks/target/release/liblambdaworks.a pkg/lambdaworks/lib
	@go build ./...

fmt:
	gofmt -w pkg

check_fmt:
	./check_fmt.sh
