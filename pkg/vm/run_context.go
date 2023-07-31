package vm

import (
	"math"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

// RunContext containts the register states of the
// Cairo VM.
type RunContext struct {
	pc memory.Relocatable
	ap uint
	fp uint
}

func (run_context RunContext) GetAp() memory.Relocatable {
	return memory.NewRelocatable(1, run_context.ap)
}

func (run_context RunContext) GetFP() memory.Relocatable {
	return memory.NewRelocatable(1, run_context.fp)
}

func (run_context RunContext) get_pc() memory.Relocatable {
	return run_context.pc
}

func (run_context RunContext) ComputeDstAddr(instruction Instruction) (memory.Relocatable, error) {
	var base_addr memory.Relocatable
	switch instruction.DstReg {
	case AP:
		base_addr = run_context.GetAp()
	case FP:
		base_addr = run_context.GetFP()
	}

	if instruction.OffOp0 < 0 {
		return base_addr.SubRelocatable(uint(math.Abs(float64(instruction.OffDst))))
	} else {
		return base_addr.AddRelocatable(uint(instruction.OffDst))
	}

}

func (run_context RunContext) ComputeOp0Addr(instruction Instruction) (memory.Relocatable, error) {
	var base_addr memory.Relocatable
	switch instruction.Op0Reg {
	case AP:
		base_addr = run_context.GetAp()
	case FP:
		base_addr = run_context.GetFP()
	}

	if instruction.OffOp1 < 0 {
		return base_addr.SubRelocatable(uint(math.Abs(float64(instruction.OffOp0))))
	} else {
		return base_addr.AddRelocatable(uint(instruction.OffOp0))
	}
}

func (run_context RunContext) ComputeOp1Addr(instruction Instruction, op0 memory.MaybeRelocatable) (memory.Relocatable, error) {
	var base_addr memory.Relocatable

	switch instruction.Op1Addr {
	case Op1SrcFP:
		base_addr = run_context.GetFP()
	case Op1SrcAP:
		base_addr = run_context.GetAp()
	case Op1SrcImm:
		if instruction.OffOp1 == 1 {
			base_addr = run_context.get_pc()
		} else {
			base_addr = memory.NewRelocatable(-1, 0)
			return base_addr, VirtualMachineError{Msg: "UnknownOp0"}
		}
		// Todo:check case op0
	}
	if instruction.OffOp1 < 0 {
		return base_addr.SubRelocatable(uint(math.Abs(float64(instruction.OffOp1))))
	} else {
		return base_addr.AddRelocatable(uint(instruction.OffOp1))
	}
}
