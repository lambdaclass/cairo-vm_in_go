package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

const SIGNATURE_BUILTIN_NAME = "signature"

// Notice changing this to any other number breaks the code
const CELLS_PER_INSTANCE = 2

type Signature struct {
	r lambdaworks.Felt
	s lambdaworks.Felt
}

type SignatureBuiltinRunner struct {
	base       memory.Relocatable
	included   bool
	signatures map[memory.Relocatable](Signature)
}

func (signatureRunner *SignatureBuiltinRunner) Base() memory.Relocatable {
	return signatureRunner.base
}

func (signatureRunner *SignatureBuiltinRunner) Name() string {
	return SIGNATURE_BUILTIN_NAME
}

func (signatureRunner *SignatureBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	signatureRunner.base = segments.AddSegment()
}

func (signatureRunner *SignatureBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if signatureRunner.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(signatureRunner.base)}
	} else {
		return nil
	}
}

func (signatureRunner *SignatureBuiltinRunner) DeduceMemoryCell(addr memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	return nil, nil
}

func ValidationRuleSignature(mem *memory.Memory, address memory.Relocatable, signatureBuiltin *SignatureBuiltinRunner) ([]memory.Relocatable, error) {
	cell_index := address.Offset % CELLS_PER_INSTANCE
	var pub_key_address, message_addr memory.Relocatable

	if cell_index == 0 {
		pub_key_address = address
		message_addr = address + 1
		// This should be 1, since CELLS_PER_INSTANCE is 2.
	} else {
		pub_key_address = address - 1
		message_addr = address
	}

	pub_key, err1 := mem.GetFelt(pub_key_address)
	message, err2 := mem.GetFelt(message_addr)
	signature = signatureBuiltin.signatures[pub_key_address]

	// Here should go the verify with some FFI

}
