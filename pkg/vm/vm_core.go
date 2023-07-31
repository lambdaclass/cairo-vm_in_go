package vm

import (
	"fmt"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	runContext  RunContext
	currentStep uint
	segments    memory.MemorySegmentManager
}

type Operands struct {
	dst memory.MaybeRelocatable
	res memory.MaybeRelocatable
	op0 memory.MaybeRelocatable
	op1 memory.MaybeRelocatable
}

type OperandsAddresses struct {
	dst_addr memory.Relocatable
	op0_addr memory.Relocatable
	op1_addr memory.Relocatable
}

func (vm VirtualMachine) compute_operands(instruction Instruction) {}

func (vm VirtualMachine) run_instrucion(instruction Instruction) {
	fmt.Println("hello from instruction")
}
