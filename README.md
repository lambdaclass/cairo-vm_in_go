# cairo-vm.go

This is a work in progress implementation of the [Cairo VM](https://github.com/lambdaclass/cairo-vm) in `Go`. The reasons for doing this include:

- Having a diversity of implementations helps find bugs and make the whole ecosystem more resilient.
- It's a good opportunity to extensively document the VM in general, as currently the documentation on its internals is very scarce and mostly lives on the minds of a few people.

## Other docs

- [Project layout](docs/layout.md)
- [Rust/lambdaworks integration](docs/rust-integration.md)

## Installation

Go needs to be installed. For mac computers, run

```shell
brew install go
```

We also use [pyenv](https://github.com/pyenv/pyenv) to install testing-related dependencies

## Compiling, running, testing

To compile, run:

```shell
make build
```

To run the main example file, run:

```shell
make run
```

Before running the tests, install the testing dependencies:

```
make deps
```

To run all tests, activate the venv created by make deps and run the test target:

```shell
. cairo-vm-env/bin/activate
make test
```

## Project Guidelines

- PRs addressing performance are forbidden. We are currently concerned with making it work without bugs and nothing more.
- All PRs must contain tests. Code coverage has to be above 98%.
- To check for security and other types of bugs, the code will be fuzzed extensively.
- PRs must be accompanied by its corresponding documentation. A book will be written documenting the entire inner workings of it, so anyone can dive in to a Cairo VM codebase and follow it along.

# Roadmap

## First milestone: Fibonacci/Factorial

What we have:

- Parsing of `json` programs
- Decoding of instructions
- Memory relocation

What we need to finish this milestone:
- Instruction execution.
- Writing of the trace and memory into files with the correct format.
- Make the fibonacci test pass, comparing our own trace with the Rust VM one, making sure they match.
- Make the factorial test pass, same as above.


## Cairo 0/Cairo 1

The above will work for Cairo 0 programs. Cairo 1 has the following extra issues to address:

- There is no `Json` representation of a Cairo 1 Program, so we can only run contracts. This means we will have to add functions to run cairo contracts (aka implement run_from_entrypoint).
- Cairo 1 contracts use the `range_check` builtin by default, so we need to implement it.

## Full VM implementation

To have a full implementation, we will need the following:

- Builtins. Add the `BuiltinRunner` logic, then implement each builtin:
    - `Range check (RC)`
    - `Output`
    - `Bitwise`
    - `Ec_op`
    - `Pedersen`
    - `Keccak`
    - `Poseidon`
    - `Signature`
    - `Segment Arena`
- Memory layouts. This is related to builtins but we will do it after implementing them.
- Hints. Add the `HintProcessor` logic, then implement each hint. Hints need to be documented extensively, implementing them is super easy since it's just porting code; what's not so clear is what they are used for. Having explanations and examples for each is important. List of hints below:
    - Parsing of references https://github.com/lambdaclass/cairo-vm/tree/main/docs/references_parsing
    - `CommonLib` https://github.com/starkware-libs/cairo-lang/tree/master/src/starkware/cairo/common
    - `Secp`
    - `Vrf`
    - `BigInt`
    - `Blake2`
    - `DictManager`
    - `EcRecover`
    - `Field Arithmetic`
    - `Garaga`
    - `set_add`
    - `sha256 utils`
    - `ECDSA verification`
    - `uint384` and `uint384 extension`
    - `uint512 utils`
    - `Cairo 1` hints.
- Proof mode. It's important to explain in detail what this is when we do it. It's one of the most obscure parts of the VM in my experience.
- Air Public inputs. (Tied to Proof-mode)
- Temporary segments.
- Program tests from Cairo VM in Rust.
- Fuzzing/Differential fuzzing.

# Documentation

## High Level Overview

The Cairo virtual machine is meant to be used in the context of STARK validity proofs. What this means is that the point of Cairo is not just to execute some code and get a result, but to *prove* to someone else that said execution was done correctly, without them having to re-execute the entire thing. The rough flow for it looks like this:

- A user writes a Cairo program.
- The program is compiled into Cairo's VM bytecode.
- The VM executes said code and provides a *trace* of execution, i.e. a record of the state of the machine and its memory *at every step of the computation*.
- This trace is passed on to a STARK prover, which creates a cryptographic proof from it, attesting to the correct execution of the program.
- The proof is passed to a verifier, who checks that the proof is valid in a fraction of a second, without re-executing.

The main three components of this flow are:

- A Cairo compiler to turn a program written in the [Cairo programming language](https://www.cairo-lang.org/) into bytecode.
- A Cairo VM to then execute it and generate a trace.
- [A STARK prover and verifier](https://github.com/lambdaclass/starknet_stack_prover_lambdaworks) so one party can prove correct execution, while another can verify it.

While this repo is only concerned with the second component, it's important to keep in mind the other two; especially important are the prover and verifier that this VM feeds its trace to, as a lot of its design decisions come from them. This virtual machine is designed to make proving and verifying both feasible and fast, and that makes it quite different from most other VMs you are probably used to.

## Basic VM flow

Our virtual machine has a very simple flow:

- Take a compiled cairo program as input. You can check out an example program [here](https://github.com/lambdaclass/cairo_vm.go/blob/main/cairo_programs/fibonacci.cairo), and its corresponding compiled version [here](https://github.com/lambdaclass/cairo_vm.go/blob/main/cairo_programs/fibonacci.json).
- Run the bytecode from the compiled program, doing the usual `fetch->decode->execute` loop, running until program termination.
- On every step of the execution, record the values of each register.
- Take the register values and memory at every step and write them to a file, called the `execution trace`.

Barring some simplifications we made, this is all the Cairo VM does. The two main things that stand out as radically different are the memory model and the use of `Field Elements` to perform arithmetic. Below we go into more detail on each step, and in the process explain the ommisions we made.

## Architecture

The Cairo virtual machine uses a Von Neumann architecture with a Non-deterministic read-only memory. What this means, roughly, is that memory is immutable after you've written to it (i.e. you can only write to it once); this is to make the STARK proving easier, but we won't go into that here.

### Memory Segments and Relocation

The process of memory allocation in a contiguous write-once memory region can get pretty complicated. Imagine you want to have a regular call stack, with a stack pointer pointing to the top of it and allocation and deallocation of stack frames and local variables happening throughout execution. Because memory is immutable, this cannot be done the usual way; once you allocate a new stack frame that memory is set, it can't be reused for another one later on.

Because of this, memory in Cairo is divided into `segments`. This is just a way of organizing memory more conveniently for this write-once model. Each segment is nothing more than a contiguous memory region. Segments are identified by an `index`, an integer value that uniquely identifies them.

Memory `cells` (i.e. values in memory) are identified by the index of the segment they belong to and an `offset` into said segment. Thus, the memory cell `{2,0}` is the first cell of segment number `2`.

Even though this segment model is extremely convenient for the VM's execution, the STARK prover needs to have the memory as just one contiguous region. Because of this, once execution of a Cairo program finishes, all the memory segments are collapsed into one; this process is called `Relocation`. We will go into more detail on all of this below.

### Registers

There are only three registers in the Cairo VM:

- The program counter `pc`, which points to the next instruction to be executed.
- The allocation pointer `ap`, pointing to the next unused memory cell.
- The frame pointer `fp`, pointing to the base of the current stack frame. When a new function is called, `fp` is set to the current `ap`. When the function returns, `fp` goes back to its previous value. The VM creates new segments whenever dynamic allocation is needed, so for example the cairo analog to a Rust `Vec` will have its own segment. Relocation at the end meshes everything together.

### Instruction Decoding/Execution

TODO: explain the components of an instruction (`dst_reg`, `op0_reg`, etc), what each one is used for and how they're encoded/decoded.

### Felts

Felts, or Field Elements, are cairo's basic integer type. Every variable in a cairo vm that is not a pointer is a felt. From our point of view we could say a felt in cairo is an unsigned integer in the range [0, CAIRO_PRIME). This means that all operations are done modulo CAIRO_PRIME. The CAIRO_PRIME is 0x800000000000011000000000000000000000000000000000000000000000001, which means felts can be quite big (up to 252 bits), luckily, we have the [Lambdaworks](https://github.com/lambdaclass/lambdaworks) library to help with handling these big integer values and providing fast and efficient modular arithmetic.

### More on memory

The cairo memory is made up of contiguous segments of variable length identified by their index. The first segment (index 0) is the program segment, which stores the instructions of a cairo program. The following segment (index 1) is the execution segment, which holds the values that are created along the execution of the vm, for example, when we call a function, a pointer to the next instruction after the call instruction will be stored in the execution segment which will then be used to find the next instruction after the function returns. The following group of segments are the builtin segments, one for each builtin used by the program, and which hold values used by the builtin runners. The last group of segments are the user segments, which represent data structures created by the user, for example, when creating an array on a cairo program, that array will be represented in memory as its own segment.

An address (or pointer) in cairo is represented as a `relocatable` value, which is made up of a `segment_index` and an `offset`, the `segment_index` tells us which segment the value is stored in and the `offset` tells us how many values exist between the start of the segment and the value.

As the cairo memory can hold both felts and pointers, the basic memory unit is a `maybe_relocatable`, a variable that can be either a `relocatable` or a `felt`.

While memory is continous, some gaps may be present. These gaps can be created on purpose by the user, for example by running:

```
[ap + 1] = 2;
```

Where a gap is created at ap. But they may also be created indireclty by diverging branches, as for example one branch may declare a variable that the other branch doesn't, as memory needs to be allocated for both cases if the second case is ran then a gap is left where the variable should have been written.

#### Memory API

The memory can perform the following basic operations:

- `memory_add_segment`: Creates a new, empty segment in memory and returns a pointer to its start. Values cannot be inserted into a memory segment that hasn't been previously created.

- `memory_insert`: Inserts a `maybe_relocatable` value at an address indicated by a `relocatable` pointer. For this operation to succeed, the pointer's segment_index must be an existing segment (created using `memory_add_segment`), and there mustn't be a value stored at that address, as the memory is immutable after its been written once. If there is a value already stored at that address but it is equal to the value to be inserted then the operation will be successful.

- `memory_get`: Fetches a `maybe_relocatable` value from a memory address indicated by a `relocatable` pointer.

Other operations:

- `memory_load_data`: This is a convenience method, which takes an array of `maybe_relocatable` and inserts them contiguosuly in memory by calling `memory_insert` and advancing the pointer by one after each insertion. Returns a pointer to the next free memory slot after the inserted data.

#### Memory Relocation

During execution, the memory consists of segments of varying length, and they can be accessed by indicating their segment index, and the offset within that segment. When the run is finished, a relocation process takes place, which transforms this segmented memory into a contiguous list of values. The relocation process works as follows:

1- The size of each segment is calculated (The size is equal to the highest offset within the segment + 1, and not the amount of `maybe_relocatable` values, as there can be gaps)
2- A base is assigned to each segment by accumulating the size of the previous segment. The first segment's base is set to 1.
3- All `relocatable` values are converted into a single integer by adding their `offset` value to their segment's base calculated in the previous step

For example, if we have this memory represented by address, value pairs:

    0:0 -> 1
    0:1 -> 4
    0:2 -> 7
    1:0 -> 8
    1:1 -> 0:2
    1:4 -> 0:1
    2:0 -> 1

Step 1: Calculate segment sizes:

    0 -> 3
    1 -> 5
    2 -> 1

Step 2: Assign a base to each segment:

    0 -> 1
    1 -> 4 (1 + 3)
    2 -> 9 (4 + 5)

Step 3: Convert relocatables to integers

    1 (base[0] + 0) -> 1
    2 (base[0] + 1) -> 4
    3 (base[0] + 2) -> 7
    4 (base[1] + 0) -> 8
    5 (base[1] + 1) -> 3 (base[0] + 2)
    .... (memory gaps)
    8 (base[1] + 4) -> 2 (base[0] + 1)
    9 (base[2] + 0) -> 1

### Program parsing

Go through the main parts of a compiled program `Json` file. `data` field with instructions, identifiers, program entrypoint, etc.

### Code walkthrough/Write your own Cairo VM

Let's begin by creating the basic types and structures for our VM:

### Felt

As anyone who has ever written a cairo program will know, everything in cairo is a Felt. We can think of it as our unsigned integer. In this project, we use the `Lambdaworks` library to abstract ourselves from modular arithmetic.

TODO: Instructions on how to use Lambdaworks felt from Go

### Relocatable

This is how cairo represents pointers, they are made up of `SegmentIndex`, which segment the variable is in, and `Offset`, how many values exist between the start of a segment and the variable. We represent them like this:

```go
type Relocatable struct {
	SegmentIndex int
	Offset       uint
}
```

### MaybeRelocatable

As the cairo memory can hold both felts and relocatables, we need a data type that can represent both in order to represent a basic memory unit.
We would normally use enums or unions to represent this type, but as go lacks both, we will instead hold a non-typed inner value and rely on the api to make sure we can only create MaybeRelocatable values with either Felt or Relocatable as inner type.

```go
type MaybeRelocatable struct {
	inner any
}

// Creates a new MaybeRelocatable with an Int inner value
func NewMaybeRelocatableInt(felt uint) *MaybeRelocatable {
	return &MaybeRelocatable{inner: Int{felt}}
}

// Creates a new MaybeRelocatable with a Relocatable inner value
func NewMaybeRelocatableRelocatable(relocatable Relocatable) *MaybeRelocatable {
	return &MaybeRelocatable{inner: relocatable}
}
```

We will also add some methods that will allow us access `MaybeRelocatable` inner values:

```go
// If m is Int, returns the inner value + true, if not, returns zero + false
func (m *MaybeRelocatable) GetInt() (Int, bool) {
	int, is_type := m.inner.(Int)
	return int, is_type
}

// If m is Relocatable, returns the inner value + true, if not, returns zero + false
func (m *MaybeRelocatable) GetRelocatable() (Relocatable, bool) {
	rel, is_type := m.inner.(Relocatable)
	return rel, is_type
}
```

These will allow us to safely discern between `Felt` and `Relocatable` values later on.

#### Memory

As we previously described, the memory is made up of a series of segments of variable length, each containing a continuous sequence of `MaybeRelocatable` elements. Memory is also immutable, which means that once we have written a value into memory, it can't be changed.
There are multiple valid ways to represent this memory structure, but the simplest way to represent it is by using a map, maping a `Relocatable` address to a `MaybeRelocatable` value.
As we don't have an actual representation of segments, we have to keep track of the number of segments.

```go
type Memory struct {
	data         map[Relocatable]MaybeRelocatable
	num_segments uint
}
```

Now we can define the basic memory operations:

*Insert*

Here we need to make perform some checks to make sure that the memory remains consistent with its rules:
- We must check that insertions are performed on previously-allocated segments, by checking that the address's segment_index is lower than our segment counter
- We must check that we are not mutating memory we have previously written, by checking that the memory doesn't already contain a value at that address that is not equal to the one we are inserting

```go
func (m *Memory) Insert(addr Relocatable, val *MaybeRelocatable) error {
	// Check that insertions are preformed within the memory bounds
	if addr.segmentIndex >= int(m.num_segments) {
		return errors.New("Error: Inserting into a non allocated segment")
	}

	// Check for possible overwrites
	prev_elem, ok := m.data[addr]
	if ok && prev_elem != *val {
		return errors.New("Memory is write-once, cannot overwrite memory value")
	}

	m.data[addr] = *val

	return nil
}
```

*Get*

This is the easiest operation, as we only need to fetch the value from our map:

```go
// Gets some value stored in the memory address `addr`.
func (m *Memory) Get(addr Relocatable) (*MaybeRelocatable, error) {
	value, ok := m.data[addr]

	if !ok {
		return nil, errors.New("Memory Get: Value not found")
	}

	return &value, nil
}
```

### MemorySegmentManager

In our `Memory` implementation, it looks like we need to have segments allocated before performing any valid memory operation, but we can't do so from the `Memory` api. To do so, we need to use the `MemorySegmentManager`.
The `MemorySegmentManager` is in charge of creating new segments and calculating their size during the relocation process, it has the following structure:

```go
type MemorySegmentManager struct {
	segmentSizes map[uint]uint
	Memory       Memory
}
```

And the following methods:

*Add Segment*

As we are using a map, we dont have to allocate memory for the new segment, so we only have to raise our segment counter and return the first address of the new segment:

```go
func (m *MemorySegmentManager) AddSegment() Relocatable {
	ptr := Relocatable{int(m.Memory.num_segments), 0}
	m.Memory.num_segments += 1
	return ptr
}
```

*Load Data*

This method inserts a contiguous array of values starting from a certain addres in memory, and returns the next address after the inserted values. This is useful when inserting the program's instructions in memory.
In order to perform this operation, we only need to iterate over the array, inserting each value at the address indicated by `ptr` while advancing the ptr with each iteration and then return the final ptr.

```go
func (m *MemorySegmentManager) LoadData(ptr Relocatable, data *[]MaybeRelocatable) (Relocatable, error) {
	for _, val := range *data {
		err := m.Memory.Insert(ptr, &val)
		if err != nil {
			return Relocatable{0, 0}, err
		}
		ptr.offset += 1
	}
	return ptr, nil
}
```

### RunContext

The RunContext keeps track of the vm's registers. Cairo VM only has 3 registers:

- The program counter `Pc`, which points to the next instruction to be executed.
- The allocation pointer `Ap`, pointing to the next unused memory cell.
- The frame pointer `Fp`, pointing to the base of the current stack frame. When a new function is called, `Fp` is set to the current `Ap` value. When the function returns, `Fp` goes back to its previous value.

We can represent it like this:

```go
type RunContext struct {
	Pc memory.Relocatable
	Ap memory.Relocatable
	Fp memory.Relocatable
}
```

### VirtualMachine
With all of these types and structures defined, we can build our VM:

```go
type VirtualMachine struct {
	runContext  RunContext
	currentStep uint
	segments    memory.MemorySegmentManager
}
```

To begin coding the basic execution functionality of our VM, we only need these basic fields, we will be adding more fields as we dive deeper into this guide.

### Builtins

TODO

### Hints

TODO
