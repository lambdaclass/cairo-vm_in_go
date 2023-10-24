package hints

import (
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
