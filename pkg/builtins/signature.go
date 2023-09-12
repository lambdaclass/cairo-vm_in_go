package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/starknet_crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const SIGNATURE_BUILTIN_NAME = "signature"

// Notice changing this to any other number breaks the code
const SIGNATURE_CELLS_PER_INSTANCE = 2

type Signature struct {
	R lambdaworks.Felt
	S lambdaworks.Felt
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

func SignatureVerificationError() error {
	return errors.New("Signature is not valid")
}

func (r *SignatureBuiltinRunner) Include(include bool) {
	r.included = include
}

func ValidationRuleSignature(mem *memory.Memory, address memory.Relocatable, signatureBuiltin *SignatureBuiltinRunner) ([]memory.Relocatable, error) {
	cell_index := address.Offset % SIGNATURE_CELLS_PER_INSTANCE
	var pub_key_address, message_addr memory.Relocatable

	if cell_index == 0 {
		pub_key_address = address
		message_addr = address.AddUint(1)
		// This should be 1, since CELLS_PER_INSTANCE is 2.
	} else {
		pub_key_address, _ = address.SubUint(1)
		message_addr = address
	}

	pub_key, get_pubkey_error := mem.GetFelt(pub_key_address)
	message, get_message_error := mem.GetFelt(message_addr)
	signature, found_signature := signatureBuiltin.signatures[pub_key_address]

	if !found_signature || get_pubkey_error != nil || get_message_error != nil {
		return nil, SignatureVerificationError()
	}

	if starknet_crypto.VerifySignature(pub_key, message, signature.R, signature.S) {
		return []memory.Relocatable{}, nil
	} else {
		return nil, SignatureVerificationError()
	}
}

func NewSignatureBuiltinRunner() *SignatureBuiltinRunner {
	return &SignatureBuiltinRunner{}
}

func (r *SignatureBuiltinRunner) AddValidationRule(mem *memory.Memory) {
	mem.AddValidationRule(uint(r.base.SegmentIndex), RangeCheckValidationRule)
}

// Helper function to AddSignature
func AddSignature(
	signatureBuiltin *SignatureBuiltinRunner,
	address memory.Relocatable,
	signature Signature,
) {
	signatureBuiltin.signatures[address] = signature
}
