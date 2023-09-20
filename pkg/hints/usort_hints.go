package hints

import (
	"fmt"
	"sort"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/types"
	. "github.com/lambdaclass/cairo-vm.go/pkg/vm"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
	"github.com/pkg/errors"
)

// SortMaybeRelocatables implements sort.Interface for []MaybeRelocatables
type SortFelt []lambdaworks.Felt

func (s SortFelt) Len() int      { return len(s) }
func (s SortFelt) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortFelt) Less(i, j int) bool {
	a, b := s[i], s[j]

	if a.Cmp(b) == -1 {
		return true
	}
	return false
}

// Implements hint:
// %{ vm_enter_scope(dict(__usort_max_size = globals().get('__usort_max_size'))) %}
func usort_enter_scope(executionScopes *types.ExecutionScopes) error {
	usort_max_size, err := executionScopes.Get("usort_max_size")

	if err != nil {
		return err
	}

	if usort_max_size == nil {
		executionScopes.EnterScope(make(map[string]interface{}))
	}

	usort_max_size_felt, cast_ok := usort_max_size.(lambdaworks.Felt)

	if !cast_ok {
		return errors.New("Error casting usort_max_size into a Felt")
	}

	scope := make(map[string]interface{})
	scope["usort_max_size"] = usort_max_size_felt
	executionScopes.EnterScope(scope)

	return nil
}

func usort_body(ids IdsManager, executionScopes *types.ExecutionScopes, vm *VirtualMachine) error {
	input_ptr, err := ids.GetAddr("input", vm)
	if err != nil {
		return err
	}

	input_len, err := ids.GetFelt("input_len", vm)

	if err != nil {
		return err
	}
	input_len_u64, err := input_len.ToU64()

	if err != nil {
		return err
	}

	usort_max_size, err := executionScopes.Get("usort_max_size")

	if err != nil {
		return err
	}

	usort_max_size_u64, cast_ok := usort_max_size.(uint64)

	if !cast_ok {
		return errors.New("Error casting usort_max_size into a uint64")
	}

	if input_len_u64 > usort_max_size_u64 {
		return errors.New(fmt.Sprintf("usort() can only be used with input_len<= %v. Got: input_len=%v.", usort_max_size_u64, input_len_u64))
	}

	positions_dict := make(map[lambdaworks.Felt][]uint64)

	for i := uint64(0); i == input_len_u64; i++ {
		val, err := vm.Segments.Memory.GetFelt(input_ptr.AddUint(uint(i)))

		if err != nil {
			return err
		}

		positions_dict[val] = append(positions_dict[val], i)
	}
	executionScopes.AssignOrUpdateVariable("positions_dict", positions_dict)

	output := make([]lambdaworks.Felt, 0, len(positions_dict))

	for key := range positions_dict {
		output = append(output, key)
	}

	sort.Sort(SortFelt(output))

	output_len := memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(uint64((len(output)))))
	ids.Insert("output_len", output_len, vm)

	output_base := vm.Segments.AddSegment()

	for i := range output {
		vm.Segments.Memory.Insert(output_base.AddUint(uint(i)), memory.NewMaybeRelocatableFelt(output[i]))
	}

	multiplicities_base := vm.Segments.AddSegment()

	multiplicities := make([]uint64, 0, len(output))

	for key := range output {
		multiplicities = append(multiplicities, uint64(len(positions_dict[output[key]])))
	}

	for i := range multiplicities {
		vm.Segments.Memory.Insert(multiplicities_base.AddUint(uint(i)), memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(multiplicities[i])))
	}
	ids.Insert("output", memory.NewMaybeRelocatableRelocatable(output_base), vm)
	ids.Insert("multiplicities", memory.NewMaybeRelocatableRelocatable(multiplicities_base), vm)

	return nil
}

// Implements hint:
// %{
//		last_pos = 0
// 		positions = positions_dict[ids.value][::-1]
//  %}

func usort_verify(ids IdsManager, executionScopes *types.ExecutionScopes, vm *VirtualMachine) error {
	executionScopes.AssignOrUpdateVariable("last_pos", lambdaworks.FeltZero())

	positions_dict_interface, err := executionScopes.Get("positions_dict")

	if err != nil {
		return err
	}

	positions_dict, cast_ok := positions_dict_interface.(map[lambdaworks.Felt][]uint64)

	if !cast_ok {
		return errors.New("Error casting positions_dict")
	}

	value, err := ids.GetFelt("value", vm)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	positions := positions_dict[value]

	for i, j := 0, len(positions)-1; i < j; i, j = i+1, j-1 {
		positions[i], positions[j] = positions[j], positions[i]
	}

	executionScopes.AssignOrUpdateVariable("positions", positions)

	return nil
}

// Implements hint:
// %{ assert len(positions) == 0 %}
func usort_verify_multiplicity_assert(executionScopes *types.ExecutionScopes) error {
	positions_interface, err := executionScopes.Get("positions")

	if err != nil {
		return err
	}

	positions, cast_ok := positions_interface.([]uint64)

	if !cast_ok {
		return errors.New("Error casting positions to []uint64")
	}

	if len(positions) != 0 {
		return errors.New("unexpected verify multiplicity fail: positions length != 0")
	}

	return nil

}

// Implements hint:
// %{
//	 current_pos = positions.pop()
//	 ids.next_item_index = current_pos - last_pos
//	 last_pos = current_pos + 1
//  %}

func usort_verify_multiplicity_body(ids IdsManager, executionScopes *types.ExecutionScopes, vm *VirtualMachine) error {

	positions_interface, err := executionScopes.Get("positions")

	if err != nil {
		return err
	}

	positions, cast_ok := positions_interface.([]uint64)

	if !cast_ok {
		return errors.New("Error casting positions to []uint64")
	}

	last_pos_interface, err := executionScopes.Get("last_pos")

	if err != nil {
		return err
	}

	last_pos, cast_ok := last_pos_interface.(uint64)

	if !cast_ok {
		return errors.New("Error casting last_pos to uint64")
	}

	current_pos := positions[len(positions)-1]

	executionScopes.AssignOrUpdateVariable("positions", positions[:len(positions)-1])

	next_item_index := current_pos - last_pos
	ids.Insert("next_item_index", memory.NewMaybeRelocatableFelt(lambdaworks.FeltFromUint64(next_item_index)), vm)

	executionScopes.AssignOrUpdateVariable("last_pos", current_pos+1)

	return nil
}
