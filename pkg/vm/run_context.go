package vm

import (
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// RunContext containts the register states of the
// Cairo VM.
type RunContext struct {
	Pc memory.Relocatable
	Ap memory.Relocatable
	Fp memory.Relocatable
}

func (run_context RunContext) ComputeDstAddr(instruction Instruction) (memory.Relocatable, error) {
	var base_addr memory.Relocatable
	switch instruction.DstReg {
	case AP:
		base_addr = run_context.Ap
	case FP:
		base_addr = run_context.Fp
	}

	if instruction.OffOp0 < 0 {
		return base_addr.SubUint(uint(math.Abs(float64(instruction.OffDst))))
	} else {
		return base_addr.AddUint(uint(instruction.OffDst))
	}

}

func (run_context RunContext) ComputeOp0Addr(instruction Instruction) (memory.Relocatable, error) {
	var base_addr memory.Relocatable
	switch instruction.Op0Reg {
	case AP:
		base_addr = run_context.Ap
	case FP:
		base_addr = run_context.Fp
	}

	if instruction.OffOp1 < 0 {
		return base_addr.SubUint(uint(math.Abs(float64(instruction.OffOp0))))
	} else {
		return base_addr.AddUint(uint(instruction.OffOp0))
	}
}

func (run_context RunContext) ComputeOp1Addr(instruction Instruction, op0 memory.MaybeRelocatable) (memory.Relocatable, error) {
	var base_addr memory.Relocatable

	switch instruction.Op1Addr {
	case Op1SrcFP:
		base_addr = run_context.Fp
	case Op1SrcAP:
		base_addr = run_context.Ap
	case Op1SrcImm:
		if instruction.OffOp1 == 1 {
			base_addr = run_context.Pc
		} else {
			base_addr = memory.NewRelocatable(-1, 0)
			return base_addr, &VirtualMachineError{Msg: "UnknownOp0"}
		}
		// Todo:check case op0
	}
	if instruction.OffOp1 < 0 {
		return base_addr.SubUint(uint(math.Abs(float64(instruction.OffOp1))))
	} else {
		return base_addr.AddUint(uint(instruction.OffOp1))
	}
}
