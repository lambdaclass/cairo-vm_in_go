#!/usr/bin/env sh
tests_path="cairo_programs/benchmarks"

set -e

for file in $(ls $tests_path | grep .cairo | sed -E 's/\.cairo//'); do
    echo "Running $file benchmark"

    export PATH="$(pyenv root)/shims:$PATH"

    hyperfine -w 5 \
	    -n "cairo-vm (Rust)" "cairo-vm/target/release/cairo-vm-cli $tests_path/$file.json --proof_mode --memory_file /dev/null --trace_file /dev/null --layout starknet_with_keccak" \
        -n "cairo-vm.go (Go)"  "go run cmd/cli/main.go $tests_path/$file.json --proof_mode --memory_file /dev/null --trace_file /dev/null --layout starknet_with_keccak" 
done