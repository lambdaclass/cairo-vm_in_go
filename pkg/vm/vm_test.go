package vm_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func TestUpdateFpRegular(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateRegular}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong value after fp update")
	}
}

func TestUpdateFpDstInt(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateDst}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(9)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 9}) {
		t.Errorf("Wrong value after fp update")
	}
}
func TestUpdateFpDstRelocatable(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateDst}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableRelocatable(memory.Relocatable{SegmentIndex: 0, Offset: 9})}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 9}) {
		t.Errorf("Wrong value after fp update")
	}
}

func TestUpdateFpApPlus2(t *testing.T) {
	instruction := vm.Instruction{FpUpdate: vm.FpUpdateAPPlus2}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	// Change the value of Ap offset
	vm.RunContext.Ap.Offset = 7
	err := vm.UpdateFp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateFp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Fp, memory.Relocatable{SegmentIndex: 0, Offset: 9}) {
		t.Errorf("Wrong value after fp update")
	}
}

func TestUpdateApRegular(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateRegular}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 0}) {
		t.Errorf("Wrong value after ap update")
	}
}

func TestUpdateApAdd2(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd2}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 2}) {
		t.Errorf("Wrong value after ap update")
	}
}

func TestUpdateApAdd1(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd1}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong value after ap update")
	}
}
func TestUpdateApAddWithResInt(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableInt(5)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdateAp failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Ap, memory.Relocatable{SegmentIndex: 0, Offset: 5}) {
		t.Errorf("Wrong value after ap update")
	}
}

func TestUpdateApAddWithResRel(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdateA should have failed")
	}
}

func TestUpdateApAddWithoutRes(t *testing.T) {
	instruction := vm.Instruction{ApUpdate: vm.ApUpdateAdd}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdateAp(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdateA should have failed")
	}
}

func TestUpdatePcRegularNoImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateRegular, Op1Addr: vm.Op1SrcAP}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcRegularWithImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateRegular, Op1Addr: vm.Op1SrcImm}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 2}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJumpWithRelRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJump}
	res := memory.Relocatable{SegmentIndex: 0, Offset: 5}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableRelocatable(res)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, res) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJumpWithIntRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJump}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableInt(0)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJumpWithoutRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJump}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()

	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJumpRelWithIntRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJumpRel}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableInt(5)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 5}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJumpRelWithRelRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJumpRel}
	res := memory.Relocatable{SegmentIndex: 0, Offset: 5}
	operands := vm.Operands{Res: memory.NewMaybeRelocatableRelocatable(res)}
	vm := vm.NewVirtualMachine()

	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJumpRelNoRes(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJumpRel}
	operands := vm.Operands{}
	vm := vm.NewVirtualMachine()

	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

func TestUpdatePcJnzDstIsZeroNoImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz, Op1Addr: vm.Op1SrcAP}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(0)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 1}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJnzDstIsZeroWithImm(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz, Op1Addr: vm.Op1SrcImm}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(0)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 2}) {
		t.Errorf("Wrong value after pc update")
	}
}

func TestUpdatePcJnzDstNotZeroOp1Int(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(1), Op1: *memory.NewMaybeRelocatableInt(3)}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err != nil {
		t.Errorf("UpdatePc failed with error: %s", err)
	}
	if !reflect.DeepEqual(vm.RunContext.Pc, memory.Relocatable{SegmentIndex: 0, Offset: 3}) {
		t.Errorf("Wrong value after pc update %v", vm.RunContext.Pc)
	}
}

func TestUpdatePcJnzDstNotZeroOp1Rel(t *testing.T) {
	instruction := vm.Instruction{PcUpdate: vm.PcUpdateJnz}
	operands := vm.Operands{Dst: *memory.NewMaybeRelocatableInt(1), Op1: *memory.NewMaybeRelocatableRelocatable(memory.Relocatable{})}
	vm := vm.NewVirtualMachine()
	err := vm.UpdatePc(&instruction, &operands)
	if err == nil {
		t.Errorf("UpdatePc should have failed")
	}
}

// Things we are skipping for now:
// - Initializing hint_executor and passing it to `cairo_run`
// - cairo_run_config stuff
// - Asserting expected trace values
// - Asserting memory_holes
func TestFibonacci(t *testing.T) {
	// compiledProgram := parser.Parse("../../cairo_programs/fibonacci.json")

	// TODO: Uncomment test when we have the bare minimum `CairoRun`
	// err := vm.CairoRun(compiledProgram.Data)
	// if err != nil {
	// 	t.Errorf("Program execution failed with error: %s", err)
	// }
}
