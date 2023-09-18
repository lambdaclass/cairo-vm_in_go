package builtins

import (
	"encoding/binary"
	"math/bits"

	"github.com/pkg/errors"

	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const KECCAK_CELLS_PER_INSTANCE = 16
const KECCAK_INPUT_CELLS_PER_INSTANCE = 8
const KECCAK_INPUT_BIT_LENTGH = 200
const KECCAK_INPUT_BYTES_LENTGH = 25

type KeccakBuiltinRunner struct {
	base                  Relocatable
	included              bool
	cache                 map[Relocatable]Felt
	ratio                 uint
	instancesPerComponent uint
	StopPtr               *uint
}

func NewKeccakBuiltinRunner(ratio uint) *KeccakBuiltinRunner {
	return &KeccakBuiltinRunner{ratio: ratio, cache: make(map[Relocatable]Felt), instancesPerComponent: 16}
}

func DefaultKeccakBuiltinRunner() *KeccakBuiltinRunner {
	return NewKeccakBuiltinRunner(2048)
}

const KECCAK_BUILTIN_NAME = "keccak"

func (k *KeccakBuiltinRunner) Base() Relocatable {
	return k.base
}

func (k *KeccakBuiltinRunner) Name() string {
	return KECCAK_BUILTIN_NAME
}

func (k *KeccakBuiltinRunner) InitializeSegments(segments *MemorySegmentManager) {
	k.base = segments.AddSegment()
}

func (k *KeccakBuiltinRunner) InitialStack() []MaybeRelocatable {
	if k.included {
		return []MaybeRelocatable{*NewMaybeRelocatableRelocatable(k.base)}
	} else {
		return nil
	}
}

func (k *KeccakBuiltinRunner) AddValidationRule(*Memory) {}

func (k *KeccakBuiltinRunner) DeduceMemoryCell(address Relocatable, mem *Memory) (*MaybeRelocatable, error) {
	index := address.Offset % KECCAK_CELLS_PER_INSTANCE
	if index < KECCAK_INPUT_CELLS_PER_INSTANCE {
		return nil, nil
	}

	value, ok := k.cache[address]
	if ok {
		return NewMaybeRelocatableFelt(value), nil
	}

	input_start_addr, _ := address.SubUint(index)
	output_start_address := input_start_addr.AddUint(KECCAK_INPUT_CELLS_PER_INSTANCE)

	input_message := make([]byte, 0, 25*KECCAK_INPUT_CELLS_PER_INSTANCE)

	// This for block creates an input message of 200 bytes. The slice is created
	// from 8 chunks of 25 bytes which are felts. To make sure nothing breaks, the
	// numbers are checked to need at most 25 bytes for their representation, if this
	// doesn't hold, an error will be returned.
	for i := uint(0); i < KECCAK_INPUT_CELLS_PER_INSTANCE; i++ {
		felt, err := mem.GetFelt(input_start_addr.AddUint(i))
		if err != nil {
			return nil, err
		}

		if !(felt.Bits() <= 200) {
			return nil, errors.New("Expected integer to be smaller than 2^200")
		}
		// Padding should already occur
		le_bytes := felt.ToLeBytes()
		input_message = append(input_message, le_bytes[:25]...)
	}

	output_message_bytes := append(make([]byte, 0, 200), input_message...)
	// Make sure output message is padded
	output_message_bytes = output_message_bytes[:cap(output_message_bytes)]

	// Type conversions are needed to use keccakF1600
	var output_message_u64 [25]uint64
	for i := 0; i < 25; i++ {
		output_message_u64[i] = binary.LittleEndian.Uint64(output_message_bytes[8*i : 8*i+8])
	}
	keccakF1600(&output_message_u64)

	// Convert back to bytes
	output_message := make([]byte, 0, 200)
	for _, num := range output_message_u64 {
		var chunk []byte
		chunk = binary.LittleEndian.AppendUint64(chunk, num)
		output_message = append(output_message, chunk...)
	}

	for i := uint(0); i < KECCAK_INPUT_CELLS_PER_INSTANCE; i++ {
		bytes := output_message[25*i : 25*i+25]
		var padded_bytes [32]byte
		copy(padded_bytes[:], bytes)
		felt := FeltFromLeBytes(&padded_bytes)
		k.cache[output_start_address.AddUint(i)] = felt
	}
	return NewMaybeRelocatableFelt(k.cache[address]), nil
}

// The following code was copied from https://github.com/golang/crypto/blob/a3485e174077e5296d3d4a43ca31d2d21b40be2c/sha3/keccakf.go
// rc stores the round constants for use in the ι step.
var rc = [24]uint64{
	0x0000000000000001,
	0x0000000000008082,
	0x800000000000808A,
	0x8000000080008000,
	0x000000000000808B,
	0x0000000080000001,
	0x8000000080008081,
	0x8000000000008009,
	0x000000000000008A,
	0x0000000000000088,
	0x0000000080008009,
	0x000000008000000A,
	0x000000008000808B,
	0x800000000000008B,
	0x8000000000008089,
	0x8000000000008003,
	0x8000000000008002,
	0x8000000000000080,
	0x000000000000800A,
	0x800000008000000A,
	0x8000000080008081,
	0x8000000000008080,
	0x0000000080000001,
	0x8000000080008008,
}

// keccakF1600 applies the Keccak permutation to a 1600b-wide
// state represented as a slice of 25 uint64s.
func keccakF1600(a *[25]uint64) {
	// Implementation translated from Keccak-inplace.c
	// in the keccak reference code.
	var t, bc0, bc1, bc2, bc3, bc4, d0, d1, d2, d3, d4 uint64

	for i := 0; i < 24; i += 4 {
		// Combines the 5 steps in each round into 2 steps.
		// Unrolls 4 rounds per loop and spreads some steps across rounds.

		// Round 1
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[6] ^ d1
		bc1 = bits.RotateLeft64(t, 44)
		t = a[12] ^ d2
		bc2 = bits.RotateLeft64(t, 43)
		t = a[18] ^ d3
		bc3 = bits.RotateLeft64(t, 21)
		t = a[24] ^ d4
		bc4 = bits.RotateLeft64(t, 14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i]
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc2 = bits.RotateLeft64(t, 3)
		t = a[16] ^ d1
		bc3 = bits.RotateLeft64(t, 45)
		t = a[22] ^ d2
		bc4 = bits.RotateLeft64(t, 61)
		t = a[3] ^ d3
		bc0 = bits.RotateLeft64(t, 28)
		t = a[9] ^ d4
		bc1 = bits.RotateLeft64(t, 20)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc4 = bits.RotateLeft64(t, 18)
		t = a[1] ^ d1
		bc0 = bits.RotateLeft64(t, 1)
		t = a[7] ^ d2
		bc1 = bits.RotateLeft64(t, 6)
		t = a[13] ^ d3
		bc2 = bits.RotateLeft64(t, 25)
		t = a[19] ^ d4
		bc3 = bits.RotateLeft64(t, 8)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc1 = bits.RotateLeft64(t, 36)
		t = a[11] ^ d1
		bc2 = bits.RotateLeft64(t, 10)
		t = a[17] ^ d2
		bc3 = bits.RotateLeft64(t, 15)
		t = a[23] ^ d3
		bc4 = bits.RotateLeft64(t, 56)
		t = a[4] ^ d4
		bc0 = bits.RotateLeft64(t, 27)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc3 = bits.RotateLeft64(t, 41)
		t = a[21] ^ d1
		bc4 = bits.RotateLeft64(t, 2)
		t = a[2] ^ d2
		bc0 = bits.RotateLeft64(t, 62)
		t = a[8] ^ d3
		bc1 = bits.RotateLeft64(t, 55)
		t = a[14] ^ d4
		bc2 = bits.RotateLeft64(t, 39)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		// Round 2
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[16] ^ d1
		bc1 = bits.RotateLeft64(t, 44)
		t = a[7] ^ d2
		bc2 = bits.RotateLeft64(t, 43)
		t = a[23] ^ d3
		bc3 = bits.RotateLeft64(t, 21)
		t = a[14] ^ d4
		bc4 = bits.RotateLeft64(t, 14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i+1]
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc2 = bits.RotateLeft64(t, 3)
		t = a[11] ^ d1
		bc3 = bits.RotateLeft64(t, 45)
		t = a[2] ^ d2
		bc4 = bits.RotateLeft64(t, 61)
		t = a[18] ^ d3
		bc0 = bits.RotateLeft64(t, 28)
		t = a[9] ^ d4
		bc1 = bits.RotateLeft64(t, 20)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc4 = bits.RotateLeft64(t, 18)
		t = a[6] ^ d1
		bc0 = bits.RotateLeft64(t, 1)
		t = a[22] ^ d2
		bc1 = bits.RotateLeft64(t, 6)
		t = a[13] ^ d3
		bc2 = bits.RotateLeft64(t, 25)
		t = a[4] ^ d4
		bc3 = bits.RotateLeft64(t, 8)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc1 = bits.RotateLeft64(t, 36)
		t = a[1] ^ d1
		bc2 = bits.RotateLeft64(t, 10)
		t = a[17] ^ d2
		bc3 = bits.RotateLeft64(t, 15)
		t = a[8] ^ d3
		bc4 = bits.RotateLeft64(t, 56)
		t = a[24] ^ d4
		bc0 = bits.RotateLeft64(t, 27)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc3 = bits.RotateLeft64(t, 41)
		t = a[21] ^ d1
		bc4 = bits.RotateLeft64(t, 2)
		t = a[12] ^ d2
		bc0 = bits.RotateLeft64(t, 62)
		t = a[3] ^ d3
		bc1 = bits.RotateLeft64(t, 55)
		t = a[19] ^ d4
		bc2 = bits.RotateLeft64(t, 39)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		// Round 3
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[11] ^ d1
		bc1 = bits.RotateLeft64(t, 44)
		t = a[22] ^ d2
		bc2 = bits.RotateLeft64(t, 43)
		t = a[8] ^ d3
		bc3 = bits.RotateLeft64(t, 21)
		t = a[19] ^ d4
		bc4 = bits.RotateLeft64(t, 14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i+2]
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc2 = bits.RotateLeft64(t, 3)
		t = a[1] ^ d1
		bc3 = bits.RotateLeft64(t, 45)
		t = a[12] ^ d2
		bc4 = bits.RotateLeft64(t, 61)
		t = a[23] ^ d3
		bc0 = bits.RotateLeft64(t, 28)
		t = a[9] ^ d4
		bc1 = bits.RotateLeft64(t, 20)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc4 = bits.RotateLeft64(t, 18)
		t = a[16] ^ d1
		bc0 = bits.RotateLeft64(t, 1)
		t = a[2] ^ d2
		bc1 = bits.RotateLeft64(t, 6)
		t = a[13] ^ d3
		bc2 = bits.RotateLeft64(t, 25)
		t = a[24] ^ d4
		bc3 = bits.RotateLeft64(t, 8)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc1 = bits.RotateLeft64(t, 36)
		t = a[6] ^ d1
		bc2 = bits.RotateLeft64(t, 10)
		t = a[17] ^ d2
		bc3 = bits.RotateLeft64(t, 15)
		t = a[3] ^ d3
		bc4 = bits.RotateLeft64(t, 56)
		t = a[14] ^ d4
		bc0 = bits.RotateLeft64(t, 27)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc3 = bits.RotateLeft64(t, 41)
		t = a[21] ^ d1
		bc4 = bits.RotateLeft64(t, 2)
		t = a[7] ^ d2
		bc0 = bits.RotateLeft64(t, 62)
		t = a[18] ^ d3
		bc1 = bits.RotateLeft64(t, 55)
		t = a[4] ^ d4
		bc2 = bits.RotateLeft64(t, 39)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		// Round 4
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[1] ^ d1
		bc1 = bits.RotateLeft64(t, 44)
		t = a[2] ^ d2
		bc2 = bits.RotateLeft64(t, 43)
		t = a[3] ^ d3
		bc3 = bits.RotateLeft64(t, 21)
		t = a[4] ^ d4
		bc4 = bits.RotateLeft64(t, 14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i+3]
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc2 = bits.RotateLeft64(t, 3)
		t = a[6] ^ d1
		bc3 = bits.RotateLeft64(t, 45)
		t = a[7] ^ d2
		bc4 = bits.RotateLeft64(t, 61)
		t = a[8] ^ d3
		bc0 = bits.RotateLeft64(t, 28)
		t = a[9] ^ d4
		bc1 = bits.RotateLeft64(t, 20)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc4 = bits.RotateLeft64(t, 18)
		t = a[11] ^ d1
		bc0 = bits.RotateLeft64(t, 1)
		t = a[12] ^ d2
		bc1 = bits.RotateLeft64(t, 6)
		t = a[13] ^ d3
		bc2 = bits.RotateLeft64(t, 25)
		t = a[14] ^ d4
		bc3 = bits.RotateLeft64(t, 8)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc1 = bits.RotateLeft64(t, 36)
		t = a[16] ^ d1
		bc2 = bits.RotateLeft64(t, 10)
		t = a[17] ^ d2
		bc3 = bits.RotateLeft64(t, 15)
		t = a[18] ^ d3
		bc4 = bits.RotateLeft64(t, 56)
		t = a[19] ^ d4
		bc0 = bits.RotateLeft64(t, 27)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc3 = bits.RotateLeft64(t, 41)
		t = a[21] ^ d1
		bc4 = bits.RotateLeft64(t, 2)
		t = a[22] ^ d2
		bc0 = bits.RotateLeft64(t, 62)
		t = a[23] ^ d3
		bc1 = bits.RotateLeft64(t, 55)
		t = a[24] ^ d4
		bc2 = bits.RotateLeft64(t, 39)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)
	}
}

func (k *KeccakBuiltinRunner) Include(include bool) {
	k.included = include
}

func (k *KeccakBuiltinRunner) Ratio() uint {
	return k.ratio
}

func (k *KeccakBuiltinRunner) CellsPerInstance() uint {
	return KECCAK_CELLS_PER_INSTANCE
}

func (k *KeccakBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	// This condition corresponds to an uninitialized ratio for the builtin, which should only
	// happen when layout is `dynamic`
	if k.Ratio() == 0 {
		// Dynamic layout has the exact number of instances it needs (up to a power of 2).
		used, err := segments.GetSegmentUsedSize(uint(k.base.SegmentIndex))
		if err != nil {
			return 0, err
		}
		instances := used / k.CellsPerInstance()
		components := utils.NextPowOf2(instances / k.instancesPerComponent)
		size := k.CellsPerInstance() * k.instancesPerComponent * components

		return size, nil
	}

	minStep := k.ratio * k.instancesPerComponent
	if currentStep < minStep {
		return 0, memory.InsufficientAllocatedCellsErrorMinStepNotReached(minStep, k.Name())
	}
	value, err := utils.SafeDiv(currentStep, k.ratio)

	if err != nil {
		return 0, errors.Errorf("error calculating builtin memory units: %s", err)
	}

	return k.CellsPerInstance() * value, nil
}

func (k *KeccakBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	used, err := segments.GetSegmentUsedSize(uint(k.base.SegmentIndex))
	if err != nil {
		return 0, 0, err
	}

	size, err := k.GetAllocatedMemoryUnits(segments, currentStep)
	if err != nil {
		return 0, 0, err
	}

	if used > size {
		return 0, 0, memory.InsufficientAllocatedCellsErrorWithBuiltinName(k.Name(), used, size)
	}

	return used, size, nil
}

func (runner *KeccakBuiltinRunner) GetRangeCheckUsage(memory *memory.Memory) (*uint, *uint) {
	return nil, nil
}

func (runner *KeccakBuiltinRunner) GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}

func (runner *KeccakBuiltinRunner) GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint {
	// The diluted cells are:
	// state - 25 rounds times 1600 elements.
	// parity - 24 rounds times 1600/5 elements times 3 auxiliaries.
	// after_theta_rho_pi - 24 rounds times 1600 elements.
	// theta_aux - 24 rounds times 1600 elements.
	// chi_iota_aux - 24 rounds times 1600 elements times 2 auxiliaries.
	// In total 25 * 1600 + 24 * 320 * 3 + 24 * 1600 + 24 * 1600 + 24 * 1600 * 2 = 216640.
	// But we actually allocate 4 virtual columns, of dimensions 64 * 1024, in which we embed the
	// real cells, and we don't free the unused ones.
	// So the real number is 4 * 64 * 1024 = 262144.
	result, err := utils.SafeDiv(262144, dilutedNBits)

	if err != nil {
		return 0
	}

	return result
}

func (runner *KeccakBuiltinRunner) GetMemoryAccesses(manager *memory.MemorySegmentManager) ([]memory.Relocatable, error) {
	segmentSize, err := manager.GetSegmentSize(uint(runner.Base().SegmentIndex))
	if err != nil {
		return []memory.Relocatable{}, err
	}

	var ret []memory.Relocatable

	var i uint
	for i = 0; i < segmentSize; i++ {
		ret = append(ret, memory.NewRelocatable(runner.Base().SegmentIndex, i))
	}

	return ret, nil
}

func (r *KeccakBuiltinRunner) FinalStack(segments *memory.MemorySegmentManager, pointer memory.Relocatable) (memory.Relocatable, error) {
	if r.included {
		if pointer.Offset == 0 {
			return memory.Relocatable{}, NewErrNoStopPointer(r.Name())
		}

		stopPointerAddr := memory.NewRelocatable(pointer.SegmentIndex, pointer.Offset-1)

		stopPointer, err := segments.Memory.GetRelocatable(stopPointerAddr)
		if err != nil {
			return memory.Relocatable{}, err
		}

		if r.Base().SegmentIndex != stopPointer.SegmentIndex {
			return memory.Relocatable{}, NewErrInvalidStopPointerIndex(r.Name(), stopPointer, r.Base())
		}

		numInstances, err := r.GetUsedInstances(segments)
		if err != nil {
			return memory.Relocatable{}, err
		}

		used := numInstances * r.CellsPerInstance()

		if stopPointer.Offset != used {
			return memory.Relocatable{}, NewErrInvalidStopPointer(r.Name(), used, stopPointer)
		}

		r.StopPtr = &stopPointer.Offset

		return stopPointerAddr, nil
	} else {
		r.StopPtr = new(uint)
		*r.StopPtr = 0
		return pointer, nil
	}
}

func (r *KeccakBuiltinRunner) GetUsedInstances(segments *memory.MemorySegmentManager) (uint, error) {
	usedCells, err := segments.GetSegmentUsedSize(uint(r.Base().SegmentIndex))
	if err != nil {
		return 0, nil
	}

	return utils.DivCeil(usedCells, r.CellsPerInstance()), nil
}
