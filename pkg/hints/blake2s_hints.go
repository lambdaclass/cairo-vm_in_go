package hints

import (
	"math"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

func Uint32SliceToMRSlice(u32Slice []uint32) []MaybeRelocatable {
	mRSlice := make([]MaybeRelocatable, 0, len(u32Slice))
	for _, u32 := range u32Slice {
		mRSlice = append(mRSlice, *NewMaybeRelocatableFelt(FeltFromUint(uint(u32))))
	}
	return mRSlice

}

func feltSliceToUint32Slice(feltSlice []Felt) ([]uint32, error) {
	uint32Slice := make([]uint32, 0, len(feltSlice))
	for _, felt := range feltSlice {
		val, err := felt.ToU32()
		if err != nil {
			return nil, err
		}
		uint32Slice = append(uint32Slice, val)
	}
	return uint32Slice, nil
}

// Returns the range from (baseAddr - baseAddrOffset) to (baseAddr - baseAddrOffset + Size)
func getUint32MemoryRange(baseAddr Relocatable, baseAddrOffset uint, size uint, segments *MemorySegmentManager) ([]uint32, error) {
	baseAddr, err := baseAddr.SubUint(baseAddrOffset)
	if err != nil {
		return nil, err
	}
	feltRange, err := segments.GetFeltRange(baseAddr, size)
	if err != nil {
		return nil, err
	}
	return feltSliceToUint32Slice(feltRange)
}

// Returns the u32 value at memory[baseAddr - baseAddrOffset]
func getU32FromMemory(baseAddr Relocatable, baseAddrOffset uint, memory *Memory) (uint32, error) {
	baseAddr, err := baseAddr.SubUint(baseAddrOffset)
	if err != nil {
		return 0, err
	}
	felt, err := memory.GetFelt(baseAddr)
	if err != nil {
		return 0, err
	}
	return felt.ToU32()
}

func blake2sCompute(ids IdsManager, vm *VirtualMachine) error {
	output, err := ids.GetRelocatable("output", vm)
	if err != nil {
		return err
	}
	h, err := getUint32MemoryRange(output, 26, 8, &vm.Segments)
	if err != nil {
		return err
	}
	message, err := getUint32MemoryRange(output, 18, 16, &vm.Segments)
	if err != nil {
		return err
	}
	t, err := getU32FromMemory(output, 2, &vm.Segments.Memory)
	if err != nil {
		return err
	}
	f, err := getU32FromMemory(output, 1, &vm.Segments.Memory)
	if err != nil {
		return err
	}

	newState := Blake2sCompress([8]uint32(h), [16]uint32(message), t, 0, f, 0)
	data := Uint32SliceToMRSlice(newState)

	_, err = vm.Segments.LoadData(output, &data)
	return err
}

func blake2sAddUint256(ids IdsManager, vm *VirtualMachine) error {
	// Fetch ids variables
	dataPtr, err := ids.GetRelocatable("data", vm)
	if err != nil {
		return err
	}
	low, err := ids.GetFelt("low", vm)
	if err != nil {
		return err
	}
	high, err := ids.GetFelt("high", vm)
	if err != nil {
		return err
	}
	// Hint logic
	const MASK = math.MaxUint32
	const B = 32
	mask := FeltFromUint(MASK)
	// First batch
	data := make([]MaybeRelocatable, 0, 4)
	for i := uint(0); i < 4; i++ {
		data = append(data, *NewMaybeRelocatableFelt(low.Shr(B * i).And(mask)))
	}
	dataPtr, err = vm.Segments.LoadData(dataPtr, &data)
	if err != nil {
		return err
	}
	// Second batch
	data = make([]MaybeRelocatable, 0, 4)
	for i := uint(0); i < 4; i++ {
		data = append(data, *NewMaybeRelocatableFelt(high.Shr(B * i).And(mask)))
	}
	_, err = vm.Segments.LoadData(dataPtr, &data)
	return err
}

func blake2sAddUint256Bigend(ids IdsManager, vm *VirtualMachine) error {
	// Fetch ids variables
	dataPtr, err := ids.GetRelocatable("data", vm)
	if err != nil {
		return err
	}
	low, err := ids.GetFelt("low", vm)
	if err != nil {
		return err
	}
	high, err := ids.GetFelt("high", vm)
	if err != nil {
		return err
	}
	// Hint logic
	const MASK = math.MaxUint32
	const B = 32
	mask := FeltFromUint(MASK)
	// First batch
	data := make([]MaybeRelocatable, 0, 4)
	for i := uint(0); i < 4; i++ {
		data = append(data, *NewMaybeRelocatableFelt(high.Shr(B * (3 - i)).And(mask)))
	}
	dataPtr, err = vm.Segments.LoadData(dataPtr, &data)
	if err != nil {
		return err
	}
	// Second batch
	data = make([]MaybeRelocatable, 0, 4)
	for i := uint(0); i < 4; i++ {
		data = append(data, *NewMaybeRelocatableFelt(low.Shr(B * (3 - i)).And(mask)))
	}
	_, err = vm.Segments.LoadData(dataPtr, &data)
	return err
}

func blake2sFinalize(ids IdsManager, vm *VirtualMachine) error {
	const N_PACKED_INSTANCES = 7
	blake2sPtrEnd, err := ids.GetRelocatable("blake2s_ptr_end", vm)
	if err != nil {
		return err
	}
	var message [16]uint32
	modifiedIv := IV()
	modifiedIv[0] = modifiedIv[0] ^ 0x01010020
	output := Blake2sCompress(modifiedIv, message, 0, 0, 0xffffffff, 0)
	padding := modifiedIv[:]
	padding = append(padding, message[:]...)
	padding = append(padding, 0, 0xffffffff)
	padding = append(padding, output[:]...)
	fullPadding := padding
	for i := 2; i < N_PACKED_INSTANCES; i++ {
		fullPadding = append(fullPadding, padding...)
	}
	data := Uint32SliceToMRSlice(fullPadding)
	_, err = vm.Segments.LoadData(blake2sPtrEnd, &data)
	return err
}

func blake2sFinalizeV3(ids IdsManager, vm *VirtualMachine) error {
	const N_PACKED_INSTANCES = 7
	blake2sPtrEnd, err := ids.GetRelocatable("blake2s_ptr_end", vm)
	if err != nil {
		return err
	}
	var message [16]uint32
	modifiedIv := IV()
	modifiedIv[0] = modifiedIv[0] ^ 0x01010020
	output := Blake2sCompress(modifiedIv, message, 0, 0, 0xffffffff, 0)
	padding := message[:]
	padding = append(padding, modifiedIv[:]...)
	padding = append(padding, 0, 0xffffffff)
	padding = append(padding, output[:]...)
	fullPadding := padding
	for i := 2; i < N_PACKED_INSTANCES; i++ {
		fullPadding = append(fullPadding, padding...)
	}
	data := Uint32SliceToMRSlice(fullPadding)
	_, err = vm.Segments.LoadData(blake2sPtrEnd, &data)
	return err
}
