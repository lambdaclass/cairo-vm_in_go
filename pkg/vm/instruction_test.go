package vm_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

func TestNonZeroHighBit(t *testing.T) {
	var _, err = vm.DecodeInstruction(0x94A7800080008000)
	if err != vm.NonZeroHighBitError {
		t.Error("Decoding should error out with NonZeroHighBitError")
	}
}

func TestInvalidOp1Reg(t *testing.T) {
	var _, err = vm.DecodeInstruction(0x294F800080008000)
	if err != vm.InvalidOp1RegError {
		t.Error("Decoding should error out with InvalidOp1RegError")
	}
}

func TestInvalidPcUpdate(t *testing.T) {
	var _, err = vm.DecodeInstruction(0x29A8800080008000)
	if err != vm.InvalidPcUpdateError {
		t.Error("Decoding should error out with InvalidPcUpdateError")
	}
}

func TestInvalidResLogic(t *testing.T) {
	var _, err = vm.DecodeInstruction(0x2968800080008000)
	if err != vm.InvalidResError {
		t.Error("Decoding should error out with InvalidResError")
	}
}

func TestInvalidOpcode(t *testing.T) {
	var _, err = vm.DecodeInstruction(0x3948800080008000)
	if err != vm.InvalidOpcodeError {
		t.Error("Decoding should error out with InvalidOpcodeError")
	}
}

func TestInvalidApUpdate(t *testing.T) {
	var _, err = vm.DecodeInstruction(0x2D48800080008000)
	if err != vm.InvalidApUpdateError {
		t.Error("Decoding should error out with InvalidApUpdateError")
	}
}

//	0|  opcode|ap_update|pc_update|res_logic|op1_src|op0_reg|dst_reg
//
// 15|14 13 12|    11 10|  9  8  7|     6  5|4  3  2|      1|      0
//
//	 |    CALL|      ADD|     JUMP|      ADD|    IMM|     FP|     FP
//	0  0  0  1      0  1   0  0  1      0  1 0  0  1       1       1
//	0001 0100 1010 0111 = 0x14A7; offx = 0
func TestDecodeFlagsCallAddJmpAddImmFpFp(t *testing.T) {
	var instruction, err = vm.DecodeInstruction(0x14A7800080008000)

	if err != nil {
		t.Errorf("Instruction decoding failed with error %s", err)
	}
	if instruction.DstReg != vm.FP {
		t.Error("Wrong Instruction Dst Register")
	}
	if instruction.Op0Reg != vm.FP {
		t.Error("Wrong Instruction Op0 Register")
	}
	if instruction.Op1Addr != vm.Op1SrcImm {
		t.Error("Wrong Instruction Op1 Address")
	}
	if instruction.ResLogic != vm.ResAdd {
		t.Error("Wrong Instruction Res")
	}
	if instruction.PcUpdate != vm.PcUpdateJump {
		t.Error("Wrong Instruction PcUpdate")
	}
	if instruction.ApUpdate != vm.ApUpdateAdd {
		t.Error("Wrong Instruction Ap Update")
	}
	if instruction.Opcode != vm.Call {
		t.Error("Wrong Instruction Opcode")
	}
	if instruction.FpUpdate != vm.FpUpdateAPPlus2 {
		t.Error("Wrong Instruction Fp Update")
	}
}

//	0|  opcode|ap_update|pc_update|res_logic|op1_src|op0_reg|dst_reg
//
// 15|14 13 12|    11 10|  9  8  7|     6  5|4  3  2|      1|      0
//
//	 |     RET|     ADD1| JUMP_REL|      MUL|     FP|     AP|     AP
//	0  0  1  0      1  0   0  1  0      1  0 0  1  0       0       0
//	0010 1001 0100 1000 = 0x2948; offx = 0
func TestDecodeFlagsRetAdd1JmpRelMulFpApAp(t *testing.T) {
	var instruction, err = vm.DecodeInstruction(0x2948800080008000)

	if err != nil {
		t.Errorf("Instruction decoding failed with error %s", err)
	}
	if instruction.DstReg != vm.AP {
		t.Error("Wrong Instruction Dst Register")
	}
	if instruction.Op0Reg != vm.AP {
		t.Error("Wrong Instruction Op0 Register")
	}
	if instruction.Op1Addr != vm.Op1SrcFP {
		t.Error("Wrong Instruction Op1 Address")
	}
	if instruction.ResLogic != vm.ResMul {
		t.Error("Wrong Instruction Res")
	}
	if instruction.PcUpdate != vm.PcUpdateJumpRel {
		t.Error("Wrong Instruction PcUpdate")
	}
	if instruction.ApUpdate != vm.ApUpdateAdd1 {
		t.Error("Wrong Instruction Ap Update")
	}
	if instruction.Opcode != vm.Ret {
		t.Error("Wrong Instruction Opcode")
	}
	if instruction.FpUpdate != vm.FpUpdateDst {
		t.Error("Wrong Instruction Fp Update")
	}
}

//	0|  opcode|ap_update|pc_update|res_logic|op1_src|op0_reg|dst_reg
//
// 15|14 13 12|    11 10|  9  8  7|     6  5|4  3  2|      1|      0
//
//	 |ASSRT_EQ|      ADD|      JNZ|      MUL|     AP|     AP|     AP
//	0  1  0  0      1  0   1  0  0      1  0 1  0  0       0       0
//	0100 1010 0101 0000 = 0x4A50; offx = 0
func TestDecodeFlagsAssertAddJnzMulApApAp(t *testing.T) {
	var instruction, err = vm.DecodeInstruction(0x4A50800080008000)

	if err != nil {
		t.Errorf("Instruction decoding failed with error %s", err)
	}
	if instruction.DstReg != vm.AP {
		t.Error("Wrong Instruction Dst Register")
	}
	if instruction.Op0Reg != vm.AP {
		t.Error("Wrong Instruction Op0 Register")
	}
	if instruction.Op1Addr != vm.Op1SrcAP {
		t.Error("Wrong Instruction Op1 Address")
	}
	if instruction.ResLogic != vm.ResMul {
		t.Error("Wrong Instruction Res")
	}
	if instruction.PcUpdate != vm.PcUpdateJnz {
		t.Error("Wrong Instruction PcUpdate")
	}
	if instruction.ApUpdate != vm.ApUpdateAdd1 {
		t.Error("Wrong Instruction Ap Update")
	}
	if instruction.Opcode != vm.AssertEq {
		t.Error("Wrong Instruction Opcode")
	}
	if instruction.FpUpdate != vm.FpUpdateRegular {
		t.Error("Wrong Instruction Fp Update")
	}
}

//	0|  opcode|ap_update|pc_update|res_logic|op1_src|op0_reg|dst_reg
//
// 15|14 13 12|    11 10|  9  8  7|     6  5|4  3  2|      1|      0
//
//	 |ASSRT_EQ|     ADD2|      JNZ|UNCONSTRD|    OP0|     AP|     AP
//	0  1  0  0      0  0   1  0  0      0  0 0  0  0       0       0
//	0100 0010 0000 0000 = 0x4200; offx = 0
func TestDecodeFlagsAssertAdd2JnzUnconOp0ApAp(t *testing.T) {
	var instruction, err = vm.DecodeInstruction(0x4200800080008000)

	if err != nil {
		t.Errorf("Instruction decoding failed with error %s", err)
	}
	if instruction.DstReg != vm.AP {
		t.Error("Wrong Instruction Dst Register")
	}
	if instruction.Op0Reg != vm.AP {
		t.Error("Wrong Instruction Op0 Register")
	}
	if instruction.Op1Addr != vm.Op1SrcOp0 {
		t.Error("Wrong Instruction Op1 Address")
	}
	if instruction.ResLogic != vm.ResUnconstrained {
		t.Error("Wrong Instruction Res")
	}
	if instruction.PcUpdate != vm.PcUpdateJnz {
		t.Error("Wrong Instruction PcUpdate")
	}
	if instruction.ApUpdate != vm.ApUpdateRegular {
		t.Error("Wrong Instruction Ap Update")
	}
	if instruction.Opcode != vm.AssertEq {
		t.Error("Wrong Instruction Opcode")
	}
	if instruction.FpUpdate != vm.FpUpdateRegular {
		t.Error("Wrong Instruction Fp Update")
	}
}

//	0|  opcode|ap_update|pc_update|res_logic|op1_src|op0_reg|dst_reg
//
// 15|14 13 12|    11 10|  9  8  7|     6  5|4  3  2|      1|      0
//
//	 |     NOP|  REGULAR|  REGULAR|      OP1|    OP0|     AP|     AP
//	0  0  0  0      0  0   0  0  0      0  0 0  0  0       0       0
//	0000 0000 0000 0000 = 0x0000; offx = 0
func TestDecodeFlagsNopReguReguOp1Op0ApAp(t *testing.T) {
	var instruction, err = vm.DecodeInstruction(0x0000800080008000)

	if err != nil {
		t.Errorf("Instruction decoding failed with error %s", err)
	}
	if instruction.DstReg != vm.AP {
		t.Error("Wrong Instruction Dst Register")
	}
	if instruction.Op0Reg != vm.AP {
		t.Error("Wrong Instruction Op0 Register")
	}
	if instruction.Op1Addr != vm.Op1SrcOp0 {
		t.Error("Wrong Instruction Op1 Address")
	}
	if instruction.ResLogic != vm.ResOp1 {
		t.Error("Wrong Instruction Res")
	}
	if instruction.PcUpdate != vm.PcUpdateRegular {
		t.Error("Wrong Instruction PcUpdate")
	}
	if instruction.ApUpdate != vm.ApUpdateRegular {
		t.Error("Wrong Instruction Ap Update")
	}
	if instruction.Opcode != vm.NOp {
		t.Error("Wrong Instruction Opcode")
	}
	if instruction.FpUpdate != vm.FpUpdateRegular {
		t.Error("Wrong Instruction Fp Update")
	}
}

//	0|  opcode|ap_update|pc_update|res_logic|op1_src|op0_reg|dst_reg
//
// 15|14 13 12|    11 10|  9  8  7|     6  5|4  3  2|      1|      0
//
//	 |     NOP|  REGULAR|  REGULAR|      OP1|    OP0|     AP|     AP
//	0  0  0  0      0  0   0  0  0      0  0 0  0  0       0       0
//	0000 0000 0000 0000 = 0x0000; offx = 0
func TestDecodeOffsetNegative(t *testing.T) {
	var instruction, err = vm.DecodeInstruction(0x0000800180007FFF)

	if err != nil {
		t.Errorf("Instruction decoding failed with error %s", err)
	}
	if instruction.OffOp0 != -1 {
		t.Error("Wrong Instruction Offset 0")
	}
	if instruction.OffOp1 != 0 {
		t.Error("Wrong Instruction Offset 1")
	}
	if instruction.OffDst != 1 {
		t.Error("Wrong Instruction Offset destination")
	}
}
