package vm

//  Structure of the 63-bit that form the first word of each instruction.
//  See Cairo whitepaper, page 32 - https://eprint.iacr.org/2021/1063.pdf.
// ┌─────────────────────────────────────────────────────────────────────────┐
// │                     off_dst (biased representation)                     │
// ├─────────────────────────────────────────────────────────────────────────┤
// │                     off_op0 (biased representation)                     │
// ├─────────────────────────────────────────────────────────────────────────┤
// │                     off_op1 (biased representation)                     │
// ├─────┬─────┬───────┬───────┬───────────┬────────┬───────────────────┬────┤
// │ dst │ op0 │  op1  │  res  │    pc     │   ap   │      opcode       │ 0  │
// │ reg │ reg │  src  │ logic │  update   │ update │                   │    │
// ├─────┼─────┼───┬───┼───┬───┼───┬───┬───┼───┬────┼────┬────┬────┬────┼────┤
// │  0  │  1  │ 2 │ 3 │ 4 │ 5 │ 6 │ 7 │ 8 │ 9 │ 10 │ 11 │ 12 │ 13 │ 14 │ 15 │
// └─────┴─────┴───┴───┴───┴───┴───┴───┴───┴───┴────┴────┴────┴────┴────┴────┘

// Instruction is the representation of the first word of each Cairo instruction.
// Some instructions spread over two words when they use an immediate value, so
// representing the first one with this struct is enougth.
type Instruction struct {
	offDst   int
	offOp0   int
	offOp1   int
	dstReg   Register
	op0Reg   Register
	op1Addr  Op1Src
	resLogic ResLogic
	pcUpdate PcUpdate
	apUpdate ApUpdate
	fpUpdate FpUpdate
	opcode   Opcode
}

// x-----------------------------x
// x----- Instruction flags -----x
// x-----------------------------x

type Register uint

const (
	AP Register = 0
	FP Register = 1
)

type Op1Src uint

const (
	Op1SrcImm Op1Src = 0
	Op1SrcAP  Op1Src = 1
	Op1SrcFP  Op1Src = 2
	Op1SrcOp0 Op1Src = 4
)

type ResLogic uint

const (
	ResOp1           ResLogic = 0
	ResAdd           ResLogic = 1
	ResMul           ResLogic = 2
	ResUnconstrained ResLogic = 3
)

type PcUpdate uint

const (
	PcUpdateRegular PcUpdate = 0
	PcUpdateJump    PcUpdate = 1
	PcUpdateJumpRel PcUpdate = 2
	PcUpdateJnz     PcUpdate = 3
)

type ApUpdate uint

const (
	ApUpdateRegular ApUpdate = 0
	ApUpdateAdd     ApUpdate = 1
	ApUpdateAdd1    ApUpdate = 2
	ApUpdateAdd2    ApUpdate = 3
)

type FpUpdate uint

const (
	FpUpdateRegular FpUpdate = 0
	FpUpdateAPPlus2 FpUpdate = 1
	FpUpdateDst     FpUpdate = 2
)

type Opcode uint

const (
	NOp      Opcode = 0
	AssertEq Opcode = 1
	Call     Opcode = 2
	Ret      Opcode = 4
)
