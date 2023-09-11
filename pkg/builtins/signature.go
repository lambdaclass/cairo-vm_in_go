package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/starknet_crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const SIGNATURE_BUILTIN_NAME = "signature"

// Notice changing this to any other number breaks the code
const CELLS_PER_INSTANCE = 2

type Signature struct {
	r lambdaworks.Felt
	s lambdaworks.Felt
}

/*
	SignatureBuiltinRunner {
		base: 0,
		included,
		ratio: instance_def.ratio,
		cells_per_instance: 2,
		n_input_cells: 2,
		_total_n_bits: 251,
		stop_ptr: None,
		instances_per_component: 1,
		signatures: Rc::new(RefCell::new(HashMap::new())),
	}
*/
type SignatureBuiltinRunner struct {
	base       memory.Relocatable
	included   bool
	signatures map[memory.Relocatable](Signature)
}

/*
	SignatureBuiltinRunner {
	    base: 0,
	    included,
	    ratio: instance_def.ratio,
	    cells_per_instance: 2,
	    n_input_cells: 2,
	    _total_n_bits: 251,
	    stop_ptr: None,
	    instances_per_component: 1,
	    signatures: Rc::new(RefCell::new(HashMap::new())),
	}
*/
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

func ValidationRuleSignature(mem *memory.Memory, address memory.Relocatable, signatureBuiltin *SignatureBuiltinRunner) ([]memory.Relocatable, error) {
	cell_index := address.Offset % CELLS_PER_INSTANCE
	var pub_key_address, message_addr memory.Relocatable

	if cell_index == 0 {
		pub_key_address = address
		message_addr = address.AddUint(1)
		// This should be 1, since CELLS_PER_INSTANCE is 2.
	} else {
		pub_key_address, _ = address.SubUint(1)
		message_addr = address
	}

	pub_key, _ := mem.GetFelt(pub_key_address)
	message, _ := mem.GetFelt(message_addr)
	signature := signatureBuiltin.signatures[pub_key_address]

	if starknet_crypto.VerifySignature(pub_key, message, signature.r, signature.s) {
		return []memory.Relocatable{}, nil
	} else {
		return nil, SignatureVerificationError()
	}
}

func NewSignatureBuiltinRunner(included bool) SignatureBuiltinRunner {
	return SignatureBuiltinRunner{
		memory.NewRelocatable(0, 0),
		included,
		map[memory.Relocatable]Signature{},
	}
}

/*
   pub(crate) fn new(instance_def: &EcdsaInstanceDef, included: bool) -> Self {
       SignatureBuiltinRunner {
           base: 0,
           included,
           ratio: instance_def.ratio,
           cells_per_instance: 2,
           n_input_cells: 2,
           _total_n_bits: 251,
           stop_ptr: None,
           instances_per_component: 1,
           signatures: Rc::new(RefCell::new(HashMap::new())),
       }
   }
*/
