package hints

import (
	"github.com/ebfe/keccak"
	"github.com/lambdaclass/cairo-vm.go/pkg/builtins"
	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"github.com/pkg/errors"
)

func unsafeKeccak(ids IdsManager, vm *VirtualMachine, scopes ExecutionScopes) error {
	// Fetch ids variable
	lengthFelt, err := ids.GetFelt("length", vm)
	if err != nil {
		return err
	}
	length, err := lengthFelt.ToU64()
	if err != nil {
		return err
	}
	data, err := ids.GetRelocatable("data", vm)
	if err != nil {
		return err
	}
	// Check __keccak_max_size if available
	keccakMaxSizeAny, err := scopes.Get("__keccak_max_size")
	if err == nil {
		keccakMaxSize, ok := keccakMaxSizeAny.(uint64)
		if ok {
			if length > keccakMaxSize {
				return errors.Errorf("unsafe_keccak() can only be used with length<=%d. Got: length=%d", keccakMaxSize, length)
			}
		}
	}
	keccakInput := make([]byte, 0)
	for byteIdx, wordIdx := 0, 0; byteIdx < int(length); byteIdx, wordIdx = byteIdx+16, wordIdx+1 {
		wordAddr := data.AddUint(uint(wordIdx))
		word, err := vm.Segments.Memory.GetFelt(wordAddr)
		if err != nil {
			return err
		}
		nBytes := int(length) - byteIdx
		if nBytes > 16 {
			nBytes = 16
		}

		if int(word.Bits()) > 8*nBytes {
			return errors.Errorf("Invalid word size: %s", word.ToHexString())
		}

		start := 32 - nBytes
		keccakInput = append(keccakInput, word.ToBeBytes()[start:]...)

	}

	hasher := keccak.New256()
	hasher.Write(keccakInput)
	resBytes := hasher.Sum(nil)

	highBytes := append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, resBytes[:16]...)
	lowBytes := append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, resBytes[16:32]...)

	high := FeltFromBeBytes((*[32]byte)(highBytes))
	low := FeltFromBeBytes((*[32]byte)(lowBytes))

	err = ids.Insert("high", NewMaybeRelocatableFelt(high), vm)
	if err != nil {
		return err
	}
	return ids.Insert("low", NewMaybeRelocatableFelt(low), vm)
}

func unsafeKeccakFinalize(ids IdsManager, vm *VirtualMachine) error {
	// Fetch ids variables
	startPtr, err := ids.GetStructFieldRelocatable("keccak_state", 0, vm)
	if err != nil {
		return err
	}
	endPtr, err := ids.GetStructFieldRelocatable("keccak_state", 1, vm)
	if err != nil {
		return err
	}

	// Hint Logic
	nElemsFelt, err := endPtr.Sub(startPtr)
	if err != nil {
		return err
	}
	nElems, err := nElemsFelt.ToU64()
	if err != nil {
		return err
	}
	inputFelts, err := vm.Segments.GetFeltRange(startPtr, uint(nElems))
	if err != nil {
		return err
	}
	inputBytes := make([]byte, 0, 16*nElems)
	for i := 0; i < int(nElems); i++ {
		inputBytes = append(inputBytes, inputFelts[i].ToBeBytes()[16:]...)
	}

	hasher := keccak.New256()
	hasher.Write(inputBytes)
	resBytes := hasher.Sum(nil)

	highBytes := append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, resBytes[:16]...)
	lowBytes := append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, resBytes[16:32]...)

	high := FeltFromBeBytes((*[32]byte)(highBytes))
	low := FeltFromBeBytes((*[32]byte)(lowBytes))

	err = ids.Insert("high", NewMaybeRelocatableFelt(high), vm)
	if err != nil {
		return err
	}
	return ids.Insert("low", NewMaybeRelocatableFelt(low), vm)
}

func compareBytesInWordNondet(ids IdsManager, vm *VirtualMachine, constants *map[string]Felt) error {
	nBytes, err := ids.GetFelt("n_bytes", vm)
	if err != nil {
		return err
	}
	bytesInWord, err := ids.GetConst("BYTES_IN_WORD", constants)
	if nBytes.Cmp(bytesInWord) == -1 {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
}

func compareKeccakFullRateInBytesNondet(ids IdsManager, vm *VirtualMachine, constants *map[string]Felt) error {
	nBytes, err := ids.GetFelt("n_bytes", vm)
	if err != nil {
		return err
	}
	bytesInWord, err := ids.GetConst("KECCAK_FULL_RATE_IN_BYTES", constants)
	if nBytes.Cmp(bytesInWord) != -1 {
		return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltOne()))
	}
	return vm.Segments.Memory.Insert(vm.RunContext.Ap, NewMaybeRelocatableFelt(FeltZero()))
}

func blockPermutation(ids IdsManager, vm *VirtualMachine, constants *map[string]Felt) error {
	const KECCAK_SIZE = 25
	keccakStateSizeFeltsFelt, err := ids.GetConst("KECCAK_STATE_SIZE_FELTS", constants)
	if err != nil {
		return err
	}
	if keccakStateSizeFeltsFelt.Cmp(FeltFromUint64(KECCAK_SIZE)) != 0 {
		return errors.New("Assertion failed: _keccak_state_size_felts == 25")
	}

	keccakPtr, err := ids.GetRelocatable("keccak_ptr", vm)
	if err != nil {
		return err
	}
	startPtr, err := keccakPtr.SubUint(KECCAK_SIZE)
	if err != nil {
		return err
	}
	inputFelt, err := vm.Segments.GetFeltRange(startPtr, KECCAK_SIZE)
	if err != nil {
		return err
	}

	var inputU64 [KECCAK_SIZE]uint64
	for i := 0; i < KECCAK_SIZE; i++ {
		val, err := inputFelt[i].ToU64()
		if err != nil {
			return err
		}
		inputU64[i] = val
	}

	builtins.KeccakF1600(&inputU64)

	output := make([]MaybeRelocatable, 0, KECCAK_SIZE)
	for i := 0; i < KECCAK_SIZE; i++ {
		output = append(output, *NewMaybeRelocatableFelt(FeltFromUint64(inputU64[i])))
	}

	_, err = vm.Segments.LoadData(keccakPtr, &output)
	return err
}
