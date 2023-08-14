package vm

import (
	"errors"
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// RunContext contains the register states of the
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

	if instruction.Off0 < 0 {
		return base_addr.SubUint(uint(math.Abs(float64(instruction.Off0))))
	} else {
		return base_addr.AddUint(uint(instruction.Off0))
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

	if instruction.Off1 < 0 {
		return base_addr.SubUint(uint(math.Abs(float64(instruction.Off1))))
	} else {
		return base_addr.AddUint(uint(instruction.Off1))
	}
}

func (run_context RunContext) ComputeOp1Addr(instruction Instruction, op0 *memory.MaybeRelocatable) (memory.Relocatable, error) {
	var base_addr memory.Relocatable

	switch instruction.Op1Addr {
	case Op1SrcFP:
		base_addr = run_context.Fp
	case Op1SrcAP:
		base_addr = run_context.Ap
	case Op1SrcImm:
		if instruction.Off2 == 1 {
			base_addr = run_context.Pc
		} else {
			base_addr = memory.NewRelocatable(0, 0)
			return memory.Relocatable{}, &VirtualMachineError{Msg: "UnknownOp0"}
		}
	case Op1SrcOp0:
		if op0 == nil {
			return memory.Relocatable{}, errors.New("Unknown Op0")
		}
		rel, is_rel := op0.GetRelocatable()
		if is_rel {
			base_addr = rel
		} else {
			return memory.Relocatable{}, errors.New("AddressNotRelocatable")
		}
	}

	if instruction.Off2 < 0 {
		return base_addr.SubUint(uint(math.Abs(float64(instruction.Off2))))
	} else {
		return base_addr.AddUint(uint(instruction.Off2))
	}
}
