package builtins

import (
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/starknet_crypto"
	"github.com/lambdaclass/cairo-vm.go/pkg/utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

const SIGNATURE_BUILTIN_NAME = "ecdsa"

// Notice changing this to any other number breaks the code
const SIGNATURE_CELLS_PER_INSTANCE = 2

type Signature struct {
	R lambdaworks.Felt
	S lambdaworks.Felt
}
type SignatureBuiltinRunner struct {
	base                  memory.Relocatable
	included              bool
	signatures            map[memory.Relocatable](Signature)
	ratio                 uint
	instancesPerComponent uint
	StopPtr               *uint
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

	if (get_pubkey_error != nil && cell_index == 1) || (get_message_error != nil && cell_index == 0) {
		return []memory.Relocatable{}, nil
	}

	if !found_signature || get_pubkey_error != nil || get_message_error != nil {
		return nil, SignatureVerificationError()
	}

	if starknet_crypto.VerifySignature(pub_key, message, signature.R, signature.S) {
		return []memory.Relocatable{}, nil
	} else {
		return nil, SignatureVerificationError()
	}
}

func NewSignatureBuiltinRunner(ratio uint) *SignatureBuiltinRunner {
	return &SignatureBuiltinRunner{signatures: map[memory.Relocatable]Signature{}, ratio: ratio, instancesPerComponent: 1}
}

func (r *SignatureBuiltinRunner) Ratio() uint {
	return r.ratio
}

func (r *SignatureBuiltinRunner) GetAllocatedMemoryUnits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	// This condition corresponds to an uninitialized ratio for the builtin, which should only
	// happen when layout is `dynamic`
	if r.Ratio() == 0 {
		// Dynamic layout has the exact number of instances it needs (up to a power of 2).
		used, err := segments.GetSegmentUsedSize(uint(r.base.SegmentIndex))
		if err != nil {
			return 0, err
		}
		instances := used / r.CellsPerInstance()
		components := utils.NextPowOf2(instances / r.instancesPerComponent)
		size := r.CellsPerInstance() * r.instancesPerComponent * components

		return size, nil
	}

	minStep := r.ratio * r.instancesPerComponent
	if currentStep < minStep {
		return 0, memory.InsufficientAllocatedCellsErrorMinStepNotReached(minStep, r.Name())
	}
	value, err := utils.SafeDiv(currentStep, r.ratio)

	if err != nil {
		return 0, errors.Errorf("error calculating builtin memory units: %s", err)
	}

	return r.CellsPerInstance() * value, nil
}

func (r *SignatureBuiltinRunner) CellsPerInstance() uint {
	return SIGNATURE_CELLS_PER_INSTANCE
}

func (r *SignatureBuiltinRunner) GetRangeCheckUsage(memory *memory.Memory) (*uint, *uint) {
	return nil, nil
}

func (r *SignatureBuiltinRunner) GetUsedCellsAndAllocatedSizes(segments *memory.MemorySegmentManager, currentStep uint) (uint, uint, error) {
	used, err := segments.GetSegmentUsedSize(uint(r.base.SegmentIndex))
	if err != nil {
		return 0, 0, err
	}

	size, err := r.GetAllocatedMemoryUnits(segments, currentStep)
	if err != nil {
		return 0, 0, err
	}

	if used > size {
		return 0, 0, memory.InsufficientAllocatedCellsErrorWithBuiltinName(r.Name(), used, size)
	}

	return used, size, nil
}

func (r *SignatureBuiltinRunner) GetUsedDilutedCheckUnits(dilutedSpacing uint, dilutedNBits uint) uint {
	return 0
}

func (r *SignatureBuiltinRunner) GetUsedPermRangeCheckLimits(segments *memory.MemorySegmentManager, currentStep uint) (uint, error) {
	return 0, nil
}

func (r *SignatureBuiltinRunner) AddValidationRule(mem *memory.Memory) {
	mem.AddValidationRule(uint(r.base.SegmentIndex), RangeCheckValidationRule)
}

// Helper function to AddSignature
func (r *SignatureBuiltinRunner) AddSignature(
	address memory.Relocatable,
	signature Signature,
) {
	r.signatures[address] = signature
}

func (runner *SignatureBuiltinRunner) GetMemoryAccesses(manager *memory.MemorySegmentManager) ([]memory.Relocatable, error) {
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

func (r *SignatureBuiltinRunner) FinalStack(segments *memory.MemorySegmentManager, pointer memory.Relocatable) (memory.Relocatable, error) {
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

func (r *SignatureBuiltinRunner) GetUsedInstances(segments *memory.MemorySegmentManager) (uint, error) {
	usedCells, err := segments.GetSegmentUsedSize(uint(r.Base().SegmentIndex))
	if err != nil {
		return 0, nil
	}

	return utils.DivCeil(usedCells, r.CellsPerInstance()), nil
}
