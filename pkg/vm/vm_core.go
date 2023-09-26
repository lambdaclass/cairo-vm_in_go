package vm

import (
	"bytes"
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const RC_OFFSET_BITS = 16

type VirtualMachineError struct {
	Msg string
}

func (e *VirtualMachineError) Error() string {
	return fmt.Sprintf(e.Msg)
}

// VirtualMachine represents the Cairo VM.
// Runs Cairo assembly and produces an execution trace.
type VirtualMachine struct {
	RunContext      RunContext
	CurrentStep     uint
	Segments        memory.MemorySegmentManager
	BuiltinRunners  []builtins.BuiltinRunner
	Trace           []TraceEntry
	RelocatedTrace  []RelocatedTraceEntry
	RelocatedMemory map[uint]lambdaworks.Felt
	RunFinished     bool
	RcLimitsMin     *int
	RcLimitsMax     *int
}

func NewVirtualMachine() *VirtualMachine {
	segments := memory.NewMemorySegmentManager()
	builtin_runners := make([]builtins.BuiltinRunner, 0, 9) // There will be at most 9 builtins
	trace := make([]TraceEntry, 0)
	relocatedTrace := make([]RelocatedTraceEntry, 0)
	return &VirtualMachine{Segments: segments, BuiltinRunners: builtin_runners, Trace: trace, RelocatedTrace: relocatedTrace}
}

func (v *VirtualMachine) Step(hintProcessor HintProcessor, hintDataMap *map[uint][]any, constants *map[string]lambdaworks.Felt, execScopes *types.ExecutionScopes) error {
	// Run Hint
	hintDatas, ok := (*hintDataMap)[v.RunContext.Pc.Offset]
	if ok {
		for i := 0; i < len(hintDatas); i++ {
			err := hintProcessor.ExecuteHint(v, &hintDatas[i], constants, execScopes)
			if err != nil {
				return err
			}
		}
	}

	// Run Instruction
	encoded_instruction, err := v.Segments.Memory.Get(v.RunContext.Pc)
	if err != nil {
		return fmt.Errorf("Failed to fetch instruction at %+v", v.RunContext.Pc)
	}

	encoded_instruction_felt, ok := encoded_instruction.GetFelt()
	if !ok {
		return errors.New("Wrong instruction encoding")
	}

	encoded_instruction_uint, err := encoded_instruction_felt.ToU64()
	if err != nil {
		return err
	}

	instruction, err := DecodeInstruction(encoded_instruction_uint)
	if err != nil {
		return err
	}

	return v.RunInstruction(&instruction)
}

func (v *VirtualMachine) RunInstruction(instruction *Instruction) error {
	operands, operandsAddresses, err := v.ComputeOperands(*instruction)
	if err != nil {
		return err
	}

	err = v.OpcodeAssertions(*instruction, operands)
	if err != nil {
		return err
	}

	v.Trace = append(v.Trace, TraceEntry{Pc: v.RunContext.Pc, Ap: v.RunContext.Ap, Fp: v.RunContext.Fp})

	v.Segments.Memory.MarkAsAccessed(operandsAddresses.DstAddr)
	v.Segments.Memory.MarkAsAccessed(operandsAddresses.Op0Addr)
	v.Segments.Memory.MarkAsAccessed(operandsAddresses.Op1Addr)

	var off0 int = instruction.Off0 + (1 << (RC_OFFSET_BITS - 1))
	var off1 int = instruction.Off1 + (1 << (RC_OFFSET_BITS - 1))
	var off2 int = instruction.Off2 + (1 << (RC_OFFSET_BITS - 1))

	if v.RcLimitsMax == nil {
		v.RcLimitsMax = new(int)
		*v.RcLimitsMax = off0
	} else {
		var value int
		value = utils.MaxInt(*v.RcLimitsMax, off0)
		value = utils.MaxInt(value, off1)
		value = utils.MaxInt(value, off2)
		*v.RcLimitsMax = value
	}

	if v.RcLimitsMin == nil {
		v.RcLimitsMin = new(int)
		*v.RcLimitsMin = off0
	} else {
		var value int
		value = utils.MinInt(*v.RcLimitsMin, off0)
		value = utils.MinInt(value, off1)
		value = utils.MinInt(value, off2)
		*v.RcLimitsMin = value
	}

	err = v.UpdateRegisters(instruction, &operands)
	if err != nil {
		return err
	}

	v.CurrentStep++
	return nil
}

// Relocates the VM's trace, turning relocatable registers to numbered ones
func (v *VirtualMachine) RelocateTrace(relocationTable *[]uint) error {
	if len(*relocationTable) < 2 {
		return errors.New("No relocation found for execution segment")
	}

	for _, entry := range v.Trace {
		v.RelocatedTrace = append(v.RelocatedTrace, RelocatedTraceEntry{
			Pc: lambdaworks.FeltFromUint64(uint64(entry.Pc.RelocateAddress(relocationTable))),
			Ap: lambdaworks.FeltFromUint64(uint64(entry.Ap.RelocateAddress(relocationTable))),
			Fp: lambdaworks.FeltFromUint64(uint64(entry.Fp.RelocateAddress(relocationTable))),
		})
	}

	return nil
}

func (v *VirtualMachine) GetRelocatedTrace() ([]RelocatedTraceEntry, error) {
	if len(v.RelocatedTrace) > 0 {
		return v.RelocatedTrace, nil
	} else {
		return nil, errors.New("Trace not relocated")
	}
}

func (v *VirtualMachine) Relocate() error {
	v.Segments.ComputeEffectiveSizes()
	if len(v.Trace) == 0 {
		return nil
	}

	relocationTable, err := v.Segments.RelocateSegments()
	// This should be unreachable
	if err != nil {
		return errors.New("ComputeEffectiveSizes called but RelocateSegments still returned error")
	}

	relocatedMemory, err := v.Segments.RelocateMemory(&relocationTable)
	if err != nil {
		return err
	}

	v.RelocateTrace(&relocationTable)
	v.RelocatedMemory = relocatedMemory
	return nil
}

// TODO: Add ExecScopes to this when it's done
func (vm *VirtualMachine) EndRun() error {
	err := vm.VerifyAutoDeductions()
	if err != nil {
		return err
	}

	vm.RunFinished = true

	// TODO
	// switch execScopes.len() {
	// case 1:
	// 	return nil
	// default:
	// 	return ErrNoScope
	// }

	return nil
}

// Makes sure that all assigned memory cells are consistent with their auto deduction rules.
func (vm *VirtualMachine) VerifyAutoDeductions() error {
	for _, builtin := range vm.BuiltinRunners {
		var index = builtin.Base()
		for relocatableAddress, value := range vm.Segments.Memory.Data {
			if relocatableAddress.SegmentIndex != index.SegmentIndex {
				continue
			}

			deducedMemoryCell, err := builtin.DeduceMemoryCell(relocatableAddress, &vm.Segments.Memory)
			if err != nil {
				return err
			}

			if deducedMemoryCell == nil {
				continue
			}

			if *deducedMemoryCell != value {
				return &VirtualMachineError{fmt.Sprintf("InconsistentAutoDeduction: %s", builtin.Name())}
			}
		}
	}

	return nil
}

type Operands struct {
	Dst memory.MaybeRelocatable
	Res *memory.MaybeRelocatable
	Op0 memory.MaybeRelocatable
	Op1 memory.MaybeRelocatable
}

type OperandsAddresses struct {
	DstAddr memory.Relocatable
	Op0Addr memory.Relocatable
	Op1Addr memory.Relocatable
}

func (vm *VirtualMachine) OpcodeAssertions(instruction Instruction, operands Operands) error {
	switch instruction.Opcode {
	case AssertEq:
		if operands.Res == nil {
			return &VirtualMachineError{"UnconstrainedResAssertEq"}
		}
		if !operands.Res.IsEqual(&operands.Dst) {
			fmt.Println("operands Res: ", operands.Res)
			fmt.Println("Operands des: ", operands.Dst)
			return &VirtualMachineError{"DiffAssertValues"}
		}
	case Call:
		new_rel := vm.RunContext.Pc.AddUint(instruction.Size())
		returnPC := memory.NewMaybeRelocatableRelocatable(new_rel)

		if !operands.Op0.IsEqual(returnPC) {
			return &VirtualMachineError{"CantWriteReturnPc"}
		}

		returnFP := vm.RunContext.Fp
		dstRelocatable, _ := operands.Dst.GetRelocatable()
		if !returnFP.IsEqual(&dstRelocatable) {
			return &VirtualMachineError{"CantWriteReturnFp"}
		}
	}

	return nil
}

func (vm *VirtualMachine) DeduceDst(instruction Instruction, res *memory.MaybeRelocatable) *memory.MaybeRelocatable {
	switch instruction.Opcode {
	case AssertEq:
		return res
	case Call:
		return memory.NewMaybeRelocatableRelocatable(vm.RunContext.Fp)

	}
	return nil
}

// Deduces the value of op0 if possible (based on dst and op1). Otherwise, returns nil.
// If res is deduced in the process returns its deduced value as well.
func (vm *VirtualMachine) DeduceOp0(instruction *Instruction, dst *memory.MaybeRelocatable, op1 *memory.MaybeRelocatable) (deduced_op0 *memory.MaybeRelocatable, deduced_res *memory.MaybeRelocatable, error error) {
	switch instruction.Opcode {
	case Call:
		deduced_op0 := vm.RunContext.Pc
		deduced_op0.Offset += instruction.Size()
		return memory.NewMaybeRelocatableRelocatable(deduced_op0), nil, nil
	case AssertEq:
		switch instruction.ResLogic {
		case ResAdd:
			if dst != nil && op1 != nil {
				deduced_op0, err := dst.Sub(*op1)
				if err != nil {
					return nil, nil, err
				}
				return &deduced_op0, dst, nil
			}
		case ResMul:
			if dst != nil && op1 != nil {
				dst_felt, dst_is_felt := dst.GetFelt()
				op1_felt, op1_is_felt := op1.GetFelt()
				if dst_is_felt && op1_is_felt && !op1_felt.IsZero() {
					return memory.NewMaybeRelocatableFelt(dst_felt.Div(op1_felt)), dst, nil

				}
			}
		}
	}
	return nil, nil, nil
}

func (vm *VirtualMachine) DeduceOp1(instruction *Instruction, dst *memory.MaybeRelocatable, op0 *memory.MaybeRelocatable) (*memory.MaybeRelocatable, *memory.MaybeRelocatable, error) {
	if instruction.Opcode == AssertEq {
		switch instruction.ResLogic {
		case ResOp1:
			return dst, dst, nil
		case ResAdd:
			if op0 != nil && dst != nil {
				dst_rel, err := dst.Sub(*op0)
				if err != nil {
					return nil, nil, err
				}
				return &dst_rel, dst, nil
			}
		case ResMul:
			if op0 != nil && dst != nil {
				dst_felt, dst_is_felt := dst.GetFelt()
				op0_felt, op0_is_felt := op0.GetFelt()
				if dst_is_felt && op0_is_felt && !op0_felt.IsZero() {
					res := memory.NewMaybeRelocatableFelt(dst_felt.Div(op0_felt))
					return res, dst, nil
				}
			}
		}
	}
	return nil, nil, nil
}

func (vm *VirtualMachine) ComputeRes(instruction Instruction, op0 memory.MaybeRelocatable, op1 memory.MaybeRelocatable) (*memory.MaybeRelocatable, error) {
	switch instruction.ResLogic {
	case ResOp1:
		return &op1, nil

	case ResAdd:
		maybe_rel, err := op0.Add(op1)
		if err != nil {
			return nil, err
		}
		return &maybe_rel, nil

	case ResMul:
		num_op0, m_type := op0.GetFelt()
		num_op1, other_type := op1.GetFelt()
		if m_type && other_type {
			result := memory.NewMaybeRelocatableFelt(num_op0.Mul(num_op1))
			return result, nil
		} else {
			return nil, errors.New("ComputeResRelocatableMul")
		}

	case ResUnconstrained:
		return nil, nil
	}
	return nil, nil
}

func (vm *VirtualMachine) ComputeOperands(instruction Instruction) (Operands, OperandsAddresses, error) {
	var res *memory.MaybeRelocatable

	dstAddr, err := vm.RunContext.ComputeDstAddr(instruction)
	if err != nil {
		return Operands{}, OperandsAddresses{}, errors.New("FailedToComputeDstAddr")
	}
	dst, _ := vm.Segments.Memory.Get(dstAddr)

	op0Addr, err := vm.RunContext.ComputeOp0Addr(instruction)
	if err != nil {
		return Operands{}, OperandsAddresses{}, fmt.Errorf("FailedToComputeOp0Addr: %s", err)
	}
	op0Op, _ := vm.Segments.Memory.Get(op0Addr)

	op1Addr, err := vm.RunContext.ComputeOp1Addr(instruction, op0Op)
	if err != nil {
		return Operands{}, OperandsAddresses{}, fmt.Errorf("FailedToComputeOp1Addr: %s", err)
	}
	op1Op, _ := vm.Segments.Memory.Get(op1Addr)

	var op0 memory.MaybeRelocatable
	if op0Op != nil {
		op0 = *op0Op
	} else {
		op0, res, err = vm.ComputeOp0Deductions(op0Addr, &instruction, dst, op1Op)
		if err != nil {
			return Operands{}, OperandsAddresses{}, err
		}
	}

	var op1 memory.MaybeRelocatable
	if op1Op != nil {
		op1 = *op1Op
	} else {
		op1, err = vm.ComputeOp1Deductions(op1Addr, &instruction, dst, op0Op, res)
		if err != nil {
			return Operands{}, OperandsAddresses{}, err
		}
	}

	if res == nil {
		res, err = vm.ComputeRes(instruction, op0, op1)

		if err != nil {
			return Operands{}, OperandsAddresses{}, err
		}
	}

	if dst == nil {
		deducedDst := vm.DeduceDst(instruction, res)
		dst = deducedDst
		if dst != nil {
			vm.Segments.Memory.Insert(dstAddr, dst)
		}
	}

	operands := Operands{
		Dst: *dst,
		Op0: op0,
		Op1: op1,
		Res: res,
	}

	operandsAddresses := OperandsAddresses{
		DstAddr: dstAddr,
		Op0Addr: op0Addr,
		Op1Addr: op1Addr,
	}
	return operands, operandsAddresses, nil
}

// Runs deductions for Op0, first runs builtin deductions, if this fails, attempts to deduce it based on dst and op1
// Also returns res if it was also deduced in the process
// Inserts the deduced operand
// Fails if Op0 was not deduced or if an error arose in the process
func (vm *VirtualMachine) ComputeOp0Deductions(op0_addr memory.Relocatable, instruction *Instruction, dst *memory.MaybeRelocatable, op1 *memory.MaybeRelocatable) (deduced_op0 memory.MaybeRelocatable, deduced_res *memory.MaybeRelocatable, err error) {
	op0, err := vm.DeduceMemoryCell(op0_addr)
	if err != nil {
		return *memory.NewMaybeRelocatableFelt(lambdaworks.FeltZero()), nil, err
	}
	if op0 == nil {
		op0, deduced_res, err = vm.DeduceOp0(instruction, dst, op1)
		if err != nil {
			return *memory.NewMaybeRelocatableFelt(lambdaworks.FeltZero()), nil, err
		}
	}
	if op0 != nil {
		vm.Segments.Memory.Insert(op0_addr, op0)
	} else {
		return *memory.NewMaybeRelocatableFelt(lambdaworks.FeltZero()), nil, errors.New("Failed to compute or deduce op0")
	}
	return *op0, deduced_res, nil
}

// Runs deductions for Op1, first runs builtin deductions, if this fails, attempts to deduce it based on dst and op0
// Also updates res if it was also deduced in the process
// Inserts the deduced operand
// Fails if Op1 was not deduced or if an error arose in the process
func (vm *VirtualMachine) ComputeOp1Deductions(op1_addr memory.Relocatable, instruction *Instruction, dst *memory.MaybeRelocatable, op0 *memory.MaybeRelocatable, res *memory.MaybeRelocatable) (memory.MaybeRelocatable, error) {
	op1, err := vm.DeduceMemoryCell(op1_addr)
	if err != nil {
		return *memory.NewMaybeRelocatableFelt(lambdaworks.FeltZero()), err
	}
	if op1 == nil {
		var deducedRes *memory.MaybeRelocatable
		op1, deducedRes, err = vm.DeduceOp1(instruction, dst, op0)
		if err != nil {
			return *memory.NewMaybeRelocatableFelt(lambdaworks.FeltZero()), err
		}
		if res == nil {
			res = deducedRes
		}
	}
	if op1 != nil {
		vm.Segments.Memory.Insert(op1_addr, op1)
	} else {
		return *memory.NewMaybeRelocatableFelt(lambdaworks.FeltZero()), errors.New("Failed to compute or deduce op1")
	}
	return *op1, nil
}

// Updates the values of the RunContext's registers according to the executed instruction
func (vm *VirtualMachine) UpdateRegisters(instruction *Instruction, operands *Operands) error {
	if err := vm.UpdateFp(instruction, operands); err != nil {
		return err
	}
	if err := vm.UpdateAp(instruction, operands); err != nil {
		return err
	}
	return vm.UpdatePc(instruction, operands)
}

// Updates the value of PC according to the executed instruction
func (vm *VirtualMachine) UpdatePc(instruction *Instruction, operands *Operands) error {
	switch instruction.PcUpdate {
	case PcUpdateRegular:
		vm.RunContext.Pc.Offset += instruction.Size()
	case PcUpdateJump:
		if operands.Res == nil {
			return errors.New("Res.UNCONSTRAINED cannot be used with PcUpdate.JUMP")
		}
		res, ok := operands.Res.GetRelocatable()
		if !ok {
			return errors.New("An integer value as Res cannot be used with PcUpdate.JUMP")
		}
		vm.RunContext.Pc = res
	case PcUpdateJumpRel:
		if operands.Res == nil {
			return errors.New("Res.UNCONSTRAINED cannot be used with PcUpdate.JUMP_REL")
		}
		res, ok := operands.Res.GetFelt()
		if !ok {
			return errors.New("A relocatable value as Res cannot be used with PcUpdate.JUMP_REL")
		}
		new_pc, err := vm.RunContext.Pc.AddFelt(res)
		if err != nil {
			return err
		}
		vm.RunContext.Pc = new_pc
	case PcUpdateJnz:
		if operands.Dst.IsZero() {
			vm.RunContext.Pc.Offset += instruction.Size()
		} else {
			new_pc, err := vm.RunContext.Pc.AddMaybeRelocatable(operands.Op1)
			if err != nil {
				return err
			}
			vm.RunContext.Pc = new_pc
		}

	}
	return nil
}

// Updates the value of AP according to the executed instruction
func (vm *VirtualMachine) UpdateAp(instruction *Instruction, operands *Operands) error {
	switch instruction.ApUpdate {
	case ApUpdateAdd:
		if operands.Res == nil {
			return errors.New("Res.UNCONSTRAINED cannot be used with ApUpdate.ADD")
		}
		new_ap, err := vm.RunContext.Ap.AddMaybeRelocatable(*operands.Res)
		if err != nil {
			return err
		}
		vm.RunContext.Ap = new_ap
	case ApUpdateAdd1:
		vm.RunContext.Ap.Offset += 1
	case ApUpdateAdd2:
		vm.RunContext.Ap.Offset += 2
	}
	return nil
}

// Updates the value of FP according to the executed instruction
func (vm *VirtualMachine) UpdateFp(instruction *Instruction, operands *Operands) error {
	switch instruction.FpUpdate {
	case FpUpdateAPPlus2:
		vm.RunContext.Fp.Offset = vm.RunContext.Ap.Offset + 2
	case FpUpdateDst:
		rel, ok := operands.Dst.GetRelocatable()
		if ok {
			vm.RunContext.Fp = rel
		} else {
			felt, _ := operands.Dst.GetFelt()
			new_fp, err := vm.RunContext.Fp.AddFelt(felt)
			if err != nil {
				return err
			}
			vm.RunContext.Fp = new_fp
		}
	}
	return nil
}

// Applies the corresponding builtin's deduction rules if addr's segment index corresponds to a builtin segment
// Returns nil if there is no deduction for the address
func (vm *VirtualMachine) DeduceMemoryCell(addr memory.Relocatable) (*memory.MaybeRelocatable, error) {
	if addr.SegmentIndex < 0 {
		return nil, nil
	}
	for i := range vm.BuiltinRunners {
		if vm.BuiltinRunners[i].Base().SegmentIndex == addr.SegmentIndex {
			return vm.BuiltinRunners[i].DeduceMemoryCell(addr, &vm.Segments.Memory)
		}
	}
	return nil, nil
}

// Write the values hosted in the output builtin's segment.
// Does nothing if the output builtin is not present in the program.
func (vm *VirtualMachine) WriteOutput(writer *bytes.Buffer) {
	for _, builtin := range vm.BuiltinRunners {
		if builtin.Name() == builtins.OUTPUT_BUILTIN_NAME {
			segmentUsedSizes := vm.Segments.ComputeEffectiveSizes()
			segmentIndex := builtin.Base().SegmentIndex
			outputSegmentSize := segmentUsedSizes[uint(segmentIndex)]

			for i := 0; i < int(outputSegmentSize); i++ {
				addr := memory.NewRelocatable(segmentIndex, uint(i))
				formattedValue, err := vm.Segments.Memory.Get(addr)
				if err != nil {
					writer.WriteString("<missing>\n")
				} else {
					writer.WriteString(formattedValue.ToString())
					writer.WriteString("\n")
				}
			}
			break
		}
	}
}

func (vm *VirtualMachine) GetBuiltinRunner(builtinName string) (*builtins.BuiltinRunner, error) {

	for _, builtin := range vm.BuiltinRunners {
		if builtin.Name() == builtinName {
			return &builtin, nil
		}
	}
	return nil, &VirtualMachineError{"BuiltinNotFound"}
}

func (vm *VirtualMachine) GetRangeCheckBound() (lambdaworks.Felt, error) {
	builtin, err := vm.GetBuiltinRunner("range_check")
	if err != nil {
		return lambdaworks.FeltZero(), err
	}

	rcBuiltin, ok := (*builtin).(*builtins.RangeCheckBuiltinRunner)
	if !ok {
		return lambdaworks.FeltZero(), errors.New("could not cast to RangeCheckBuiltinRunner")
	}

	return rcBuiltin.Bound(), nil
}
