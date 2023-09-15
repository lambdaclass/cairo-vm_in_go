.PHONY: deps deps-macos run test coverage build fmt check_fmt clean clean_files build_cairo_vm_cli compare_trace_memory compare_trace \
 compare_memory demo_fibonacci demo_factorial compare_proof_trace_memory compare_proof_trace compare_proof_memory $(CAIRO_VM_CLI) \

CAIRO_VM_CLI:=cairo-vm/target/release/cairo-vm-cli

$(CAIRO_VM_CLI):
	git clone --depth 1 -b v0.8.5 https://github.com/lambdaclass/cairo-vm
	cd cairo-vm; cargo b --release --bin cairo-vm-cli

# Create proof mode programs. 
# NOTE: This is super flaky, DO NOT move this section below the non proof mode one
# or things will go wrong.

TEST_PROOF_DIR=cairo_programs/proof_programs
TEST_PROOF_FILES:=$(wildcard $(TEST_PROOF_DIR)/*.cairo)
COMPILED_PROOF_TESTS:=$(patsubst $(TEST_PROOF_DIR)/%.cairo, $(TEST_PROOF_DIR)/%.json, $(TEST_PROOF_FILES))

CAIRO_RS_PROOF_MEM:=$(patsubst $(TEST_PROOF_DIR)/%.json, $(TEST_PROOF_DIR)/%.rs.memory, $(COMPILED_PROOF_TESTS))
CAIRO_RS_PROOF_TRACE:=$(patsubst $(TEST_PROOF_DIR)/%.json, $(TEST_PROOF_DIR)/%.rs.trace, $(COMPILED_PROOF_TESTS))

CAIRO_GO_PROOF_MEM:=$(patsubst $(TEST_PROOF_DIR)/%.json, $(TEST_PROOF_DIR)/%.go.memory, $(COMPILED_PROOF_TESTS))
CAIRO_GO_PROOF_TRACE:=$(patsubst $(TEST_PROOF_DIR)/%.json, $(TEST_PROOF_DIR)/%.go.trace, $(COMPILED_PROOF_TESTS))

$(TEST_PROOF_DIR)/%.json: $(TEST_PROOF_DIR)/%.cairo
	cairo-compile --cairo_path="$(TEST_PROOF_DIR)" $< --output $@ --proof_mode

$(TEST_PROOF_DIR)/%.rs.trace $(TEST_PROOF_DIR)/%.rs.memory: $(TEST_PROOF_DIR)/%.json $(CAIRO_VM_CLI)
	$(CAIRO_VM_CLI) --layout all_cairo $< --trace_file $(@D)/$(*F).rs.trace --memory_file $(@D)/$(*F).rs.memory --proof_mode

$(TEST_PROOF_DIR)/%.go.trace $(TEST_PROOF_DIR)/%.go.memory: $(TEST_PROOF_DIR)/%.json
	go run cmd/cli/main.go --trace_file $(@D)/$(*F).go.trace --memory_file $(@D)/$(*F).go.memory --layout all_cairo --proof_mode $(@D)/$(*F).json

# Non Proof mode programs

TEST_DIR=cairo_programs
TEST_FILES:=$(wildcard $(TEST_DIR)/*.cairo)
COMPILED_TESTS:=$(patsubst $(TEST_DIR)/%.cairo, $(TEST_DIR)/%.json, $(TEST_FILES))

CAIRO_RS_MEM:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.rs.memory, $(COMPILED_TESTS))
CAIRO_RS_TRACE:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.rs.trace, $(COMPILED_TESTS))

CAIRO_GO_MEM:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.go.memory, $(COMPILED_TESTS))
CAIRO_GO_TRACE:=$(patsubst $(TEST_DIR)/%.json, $(TEST_DIR)/%.go.trace, $(COMPILED_TESTS))

$(TEST_DIR)/%.rs.trace $(TEST_DIR)/%.rs.memory: $(TEST_DIR)/%.json $(CAIRO_VM_CLI)
	$(CAIRO_VM_CLI) --layout all_cairo $< --trace_file $(@D)/$(*F).rs.trace --memory_file $(@D)/$(*F).rs.memory

$(TEST_DIR)/%.go.trace $(TEST_DIR)/%.go.memory: $(TEST_DIR)/%.json
	go run cmd/cli/main.go --trace_file $(@D)/$(*F).go.trace --memory_file $(@D)/$(*F).go.memory --layout all_cairo $(@D)/$(*F).json

$(TEST_DIR)/%.json: $(TEST_DIR)/%.cairo
	cairo-compile --cairo_path="$(TEST_DIR)" $< --output $@

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

test: build $(COMPILED_TESTS) $(COMPILED_PROOF_TESTS)
	@go test -v ./...

coverage: $(COMPILED_TESTS)
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...

coverage_html: coverage
	@go tool cover -html=coverage.out

build:
	@cd pkg/lambdaworks/lib/lambdaworks && cargo build --release
	@cp pkg/lambdaworks/lib/lambdaworks/target/release/liblambdaworks.a pkg/lambdaworks/lib
	@cd pkg/starknet_crypto/lib/starknet_crypto && cargo build --release
	@cp pkg/starknet_crypto/lib/starknet_crypto/target/release/libstarknet_crypto.a pkg/starknet_crypto/lib
	@go build ./...

fmt:
	gofmt -w pkg

check_fmt:
	./check_fmt.sh

clean:
	rm -f $(TEST_DIR)/*.json
	rm -f $(TEST_DIR)/*.memory
	rm -f $(TEST_DIR)/*.trace
	cd pkg/lambdaworks/lib/lambdaworks && cargo clean
	rm -f pkg/lambdaworks/lib/liblambdaworks.a
	cd pkg/starknet_crypto/lib/starknet_crypto && cargo clean
	rm -f pkg/starknet_crypto/lib/libstarknet_crypto.a
	rm -rf cairo-vm
	rm -rf cairo-vm-env

clean_files:
	rm -f $(TEST_DIR)/*.json
	rm -f $(TEST_DIR)/*.memory
	rm -f $(TEST_DIR)/*.trace

demo_fibonacci: clean_files build_cairo_vm_cli build
	@echo "Compiling fibonacci program..."
	@cairo-compile --cairo_path="$(TEST_DIR)" cairo_programs/fibonacci.cairo --output cairo_programs/fibonacci.json
	@echo "Running fibonacci program with Go implementation..."
	@go run cmd/cli/main.go cairo_programs/fibonacci.json
	@echo "Running fibonacci program with Rust implementation..."
	@$(CAIRO_VM_CLI) --layout all_cairo cairo_programs/fibonacci.json --trace_file cairo_programs/fibonacci.rs.trace --memory_file cairo_programs/fibonacci.rs.memory
	@echo "Done!"
	@echo "Comparing fibonacci trace with Rust implementation..."
	@if ! diff -q cairo_programs/fibonacci.go.trace cairo_programs/fibonacci.rs.trace; then \
		@echo "\xE2\x9D\x8E Traces for fibonacci differ"; \
		@exit 1; \
	fi
	@echo "\xE2\x9C\x85 Traces for fibonacci match!"
	@echo "Comparing fibonacci memory with Rust implementation..."
	@if ! python scripts/memory_comparator.py cairo_programs/fibonacci.go.memory cairo_programs/fibonacci.rs.memory; then \
		@echo "\xE2\x9D\x8E Memory for fibonacci differs"; \
		@exit 1; \
	fi
	@echo "\xE2\x9C\x85 Memory for fibonacci matches!"

demo_factorial: clean_files build_cairo_vm_cli build
	@echo "Compiling factorial program..."
	@cairo-compile --cairo_path="$(TEST_DIR)" cairo_programs/factorial.cairo --output cairo_programs/factorial.json
	@echo "Running factorial program with Go implementation..."
	@go run cmd/cli/main.go cairo_programs/factorial.json
	@echo "Running factorial program with Rust implementation..."
	@$(CAIRO_VM_CLI) --layout all_cairo cairo_programs/factorial.json --trace_file cairo_programs/factorial.rs.trace --memory_file cairo_programs/factorial.rs.memory
	@echo "Done!"
	@echo "Comparing factorial trace with Rust implementation..."
	@if ! diff -q cairo_programs/factorial.go.trace cairo_programs/factorial.rs.trace; then \
		@echo "\xE2\x9D\x8E Traces for factorial differ"; \
		exit 1; \
	fi
	@echo "\xE2\x9C\x85 Traces for factorial match!"
	@echo "Comparing factorial memory with Rust implementation..."
	@if ! python scripts/memory_comparator.py cairo_programs/factorial.go.memory cairo_programs/factorial.rs.memory; then \
		@echo "\xE2\x9D\x8E Memory for factorial differs"; \
		exit 1; \
	fi
	@echo "\xE2\x9C\x85 Memory for factorial matches!"

build_cairo_vm_cli: | $(CAIRO_VM_CLI)

compare_trace_memory: build_cairo_vm_cli $(CAIRO_RS_MEM) $(CAIRO_RS_TRACE) $(CAIRO_GO_MEM) $(CAIRO_GO_TRACE)
	cd scripts; sh compare_vm_state.sh trace memory

compare_trace: build_cairo_vm_cli $(CAIRO_RS_TRACE) $(CAIRO_GO_TRACE)
	cd scripts; sh compare_vm_state.sh trace

compare_memory: build_cairo_vm_cli $(CAIRO_RS_MEM) $(CAIRO_GO_MEM)
	cd scripts; sh compare_vm_state.sh memory

compare_proof_trace_memory: build_cairo_vm_cli $(CAIRO_RS_PROOF_MEM) $(CAIRO_RS_PROOF_TRACE) $(CAIRO_GO_PROOF_MEM) $(CAIRO_GO_PROOF_TRACE)
	cd scripts; sh compare_vm_state.sh trace memory proof_mode

compare_proof_trace: build_cairo_vm_cli $(CAIRO_RS_PROOF_TRACE) $(CAIRO_GO_PROOF_TRACE)
	cd scripts; sh compare_vm_state.sh trace proof_mode

compare_proof_memory: build_cairo_vm_cli $(CAIRO_RS_PROOF_MEM) $(CAIRO_GO_PROOF_MEM)
	cd scripts; sh compare_vm_state.sh memory proof_mode

