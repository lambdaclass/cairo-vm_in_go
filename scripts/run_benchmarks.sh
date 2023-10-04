#!/usr/bin/env sh
tests_path="cairo_programs/benchmarks"

set -e

for file in $(ls $tests_path | grep .cairo | sed -E 's/\.cairo//'); do
    echo "Running $file benchmark"

    export PATH="$(pyenv root)/shims:$PATH"

    hyperfine \
        -n "Cairo VM (CPython)" "PYENV_VERSION=3.9.15 cairo-run --proof_mode --memory_file /dev/null --trace_file /dev/null --layout starknet_with_keccak --program $tests_path/$file.json" \
	    -n "Cairo VM (PyPy)" "PYENV_VERSION=pypy3.9-7.3.9 cairo-run --proof_mode --memory_file /dev/null --trace_file /dev/null --layout starknet_with_keccak --program $tests_path/$file.json" \
	    -n "cairo-vm (Rust)" "cairo-vm/target/release/cairo-vm-cli $tests_path/$file.json --proof_mode --memory_file /dev/null --trace_file /dev/null --layout starknet_with_keccak" \
        -n "cairo-vm.go (Go)"  "go run cmd/cli/main.go $tests_path/$file.json --proof_mode --memory_file /dev/null --trace_file /dev/null --layout starknet_with_keccak" 
done
