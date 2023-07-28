# Lambdaworks wrapper

For now, this is a mock containing only a "number" function that always return 42. It's a
proof of concept of a go package that integrates with a rust library.

## Package structure

These are the relevant files:

```
./pkg/lambdaworks: go package wrapper so that other packages can call this
        │          one without worrying about FFI, C or rust.
        ├──lambdaworks.go: go wrapper library code. Casts C types to go types.
        └── lib: directory with the rust and C code.
            ├── lambdaworks.h: C headers representing the functions that will
            │                  be exported to go.
            └── lambdaworks: Rust package.
                ├── Cargo.Toml: Rust package definition.
                └── src/lib.rs: rust library.
```

## Compiling

The rust package can be compiled by moving to the `lib/lambdaworks` directory and executing `cargo build --release`. This produces a `lib/lambdaworks/target/liblambdaworks.a` file that will be imported by the `lambdaworks.go` as a static library, as we can see in the comment over the `Number` function.

In `Makefile` we can see that one of the steps is actually moving to the rust project, compiling, and then copying the archive file to the `lib` directory so it can be consumed by the go code.

## Next steps:

The lambdaworks math and crypto rust dependencies are already included in the `Cargo.toml`, but they are not used. The first step would be to add a wrapper that manipulates the finite field type instead of a simple integer.
