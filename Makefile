.PHONY: deps deps-macos run test build fmt check_fmt clean

TEST_DIR=cairo_programs
TEST_FILES:=$(wildcard $(TEST_DIR)/*.cairo)
COMPILED_TESTS:=$(patsubst $(TEST_DIR)/%.cairo, $(TEST_DIR)/%.json, $(TEST_FILES))

$(TEST_DIR)/%.json: $(TEST_DIR)/%.cairo
	. cairo-vm-env/bin/activate ; \
	cairo-compile --cairo_path="$(TEST_DIR):$(BENCH_DIR)" $< --output $@

# Creates a pyenv and installs cairo-lang
deps:
	pyenv install  -s 3.9.15
	PYENV_VERSION=3.9.15 python -m venv cairo-vm-env
	. cairo-vm-env/bin/activate ; \
	pip install -r requirements.txt ; \

# Creates a pyenv and installs cairo-lang
deps-macos:
	brew install gmp pyenv
	pyenv install -s 3.9.15
	PYENV_VERSION=3.9.15 python -m venv cairo-vm-env
	. cairo-vm-env/bin/activate ; \
	CFLAGS=-I/opt/homebrew/opt/gmp/include LDFLAGS=-L/opt/homebrew/opt/gmp/lib pip install -r requirements.txt ; \

run:
	@go run cmd/cli/main.go

test: $(COMPILED_TESTS)
	@go test -v ./...

build:
	@cd pkg/lambdaworks/lib/lambdaworks && cargo build --release
	@cp pkg/lambdaworks/lib/lambdaworks/target/release/liblambdaworks.a pkg/lambdaworks/lib
	@go build ./...

fmt:
	gofmt -w pkg

check_fmt:
	./check_fmt.sh

clean:
	rm cairo_programs/*.json
	rm -r cairo-vm-env

demo: $(COMPILED_TESTS)
	go run .
