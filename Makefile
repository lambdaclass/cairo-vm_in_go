.PHONY: deps deps-macos run test build fmt check_fmt clean build_cairo_vm_cli compare_trace_memory compare_trace compare_memory $(CAIRO_VM_CLI)

CAIRO_VM_CLI:=cairo-vm/target/release/cairo-vm-cli

$(CAIRO_VM_CLI):
	git clone --depth 1 -b v0.8.5 https://github.com/lambdaclass/cairo-vm
	cd cairo-vm; cargo b --release --bin cairo-vm-cli

TEST_DIR=cairo_programs
TEST_FILES:=$(wildcard $(TEST_DIR)/*.cairo)
COMPILED_TESTS:=$(patsubst $(TEST_DIR)/%.cairo, $(TEST_DIR)/%.json, $(TEST_FILES))

CAIRO_RS_MEM:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.rs.memory, $(COMPILED_TESTS))
CAIRO_RS_TRACE:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.rs.trace, $(COMPILED_TESTS))

CAIRO_GO_MEM:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.go.memory, $(COMPILED_TESTS))
CAIRO_GO_TRACE:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.go.trace, $(COMPILED_TESTS))

$(TEST_DIR)/%.rs.trace $(TEST_DIR)/%.rs.memory: $(TEST_DIR)/%.json $(CAIRO_VM_CLI)
	$(CAIRO_VM_CLI) --layout all_cairo $< --trace_file $@ --memory_file $(@D)/$(*F).rs.memory

# TODO: Uses cairo-lang as placeholder, should use cairo-vm.go
$(TEST_DIR)/%.go.trace $(TEST_DIR)/%.go.memory: $(TEST_DIR)/%.json
	cairo-run --layout starknet_with_keccak --program $< --trace_file $@ --memory_file $(@D)/$(*F).go.memory

$(TEST_DIR)/%.json: $(TEST_DIR)/%.cairo
	cairo-compile --cairo_path="$(TEST_DIR):$(BENCH_DIR)" $< --output $@

# Creates a pyenv and installs cairo-lang
deps:
	pyenv install  -s 3.9.15
	PYENV_VERSION=3.9.15 python -m venv cairo-vm-env
	. cairo-vm-env/bin/activate ; \
	pip install -r requirements.txt ; \

# Creates a pyenv and installs cairo-lang
deps-macos:
	brew install gmp
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
	rm -f $(TEST_DIR)/*.json
	rm -f $(TEST_DIR)/*.memory
	rm -f $(TEST_DIR)/*.trace
	rm -rf cairo-vm
	rm -r cairo-vm-env

build_cairo_vm_cli: | $(CAIRO_VM_CLI)

compare_trace_memory: build_cairo_vm_cli $(CAIRO_RS_MEM) $(CAIRO_RS_TRACE) $(CAIRO_GO_MEM) $(CAIRO_GO_TRACE)
	cd scripts; sh compare_vm_state.sh trace memory

compare_trace: build_cairo_vm_cli $(CAIRO_RS_TRACE) $(CAIRO_GO_TRACE)
	cd scripts; sh compare_vm_state.sh trace

compare_memory: build_cairo_vm_cli $(CAIRO_RS_MEM) $(CAIRO_GO_MEM)
	cd scripts; sh compare_vm_state.sh memory
