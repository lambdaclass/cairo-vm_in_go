package builtins

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
)

type EcOpBuiltinRunner struct {
	included           bool
	base               memory.Relocatable
	cells_per_instance uint32
	cache              map[memory.Relocatable]lambdaworks.Felt
	n_input_cells      uint32
	scalar_height      uint32
}

type EcPoint struct {
	x uint
	y uint
}

type PartialSum struct {
	x lambdaworks.Felt
	y lambdaworks.Felt
}

type DoublePoint struct {
	x lambdaworks.Felt
	y lambdaworks.Felt
}

const INPUT_CELLS_PER_EC_OP = 5
const PRIME = "0x800000000000011000000000000000000000000000000000000000000000001"

func NewEcOpBuiltinRunner(included bool) *EcOpBuiltinRunner {
	return &EcOpBuiltinRunner{
		included: included,
	}
}

func (ec *EcOpBuiltinRunner) Base() memory.Relocatable {
	return ec.base
}

func (ec *EcOpBuiltinRunner) Name() string {
	return "ec_op"
}

func (ec *EcOpBuiltinRunner) InitializeSegments(segments *memory.MemorySegmentManager) {
	ec.base = segments.AddSegment()
}

func (ec *EcOpBuiltinRunner) InitialStack() []memory.MaybeRelocatable {
	if ec.included {
		return []memory.MaybeRelocatable{*memory.NewMaybeRelocatableRelocatable(ec.base)}
	} else {
		return []memory.MaybeRelocatable{}
	}
}

func (ec *EcOpBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, segments *memory.Memory) (*memory.MaybeRelocatable, error) {
	EC_POINT_INDICES := [3]EcPoint{EcPoint{x: 0, y: 1}, EcPoint{x: 2, y: 3}, EcPoint{x: 5, y: 6}}
	OUTPUT_INDICES := EC_POINT_INDICES[2]
	alpha := lambdaworks.FeltOne()
	beta_low := lambdaworks.FeltFromHex("0x609ad26c15c915c1f4cdfcb99cee9e89")
	beta_high := lambdaworks.FeltFromHex("0x6f21413efbe40de150e596d72f7a8c5")
	beta := (beta_high.Shl(128)).Add(beta_low)

	index := address.Offset % uint(ec.cells_per_instance)

	if index != OUTPUT_INDICES.x && index != OUTPUT_INDICES.y {
		return nil, nil
	}

	instance := memory.NewRelocatable(address.SegmentIndex, address.Offset-index)
	input_cells_per_ec_op := lambdaworks.FeltFromUint64(INPUT_CELLS_PER_EC_OP)
	x_addr, err := instance.AddFelt(input_cells_per_ec_op)
	if err != nil {
		return nil, errors.New("Runner error, Expected Integer")
	}

	number := ec.cache[address]

	if !reflect.DeepEqual(number, lambdaworks.Felt{}) {
		return memory.NewMaybeRelocatableFelt(number), nil
	}

	//All input cells should be filled, and be integer values
	//If an input cell is not filled, return None

	input_cells := make([]lambdaworks.Felt, ec.n_input_cells)
	for i := 0; i < int(ec.n_input_cells); i++ {
		addr, err := segments.Get(instance.AddUint(uint(i)))
		if err != nil {
			felt, is_felt := addr.GetFelt()
			if is_felt {
				input_cells = append(input_cells, felt)
			} else {
				return nil, errors.New("Runner error, Expected Integer for input cells")
			}
		} else {
			return nil, nil
		}
	}

	for j := 0; j < len(EC_POINT_INDICES); j++ {
		x := input_cells[EC_POINT_INDICES[j].x]
		y := input_cells[EC_POINT_INDICES[j].y]
		if !ec.point_on_curve(x, y, alpha, beta) {
			return nil, errors.New("Point not in curve")
		}
	}

	prime, ok := new(big.Int).SetString(PRIME, 16)
	if !ok {
		return nil, errors.New("Could not parse prime")
	}

	alpha_big_int := big.NewInt(1)

	partial_sum := PartialSum{x: input_cells[0], y: input_cells[1]}
	double_point := DoublePoint{x: input_cells[2], y: input_cells[3]}

	result_x, result_y := ec.ec_op_impl(partial_sum, double_point, input_cells[4], alpha_big_int, prime, ec.scalar_height)

	felt_result_x := lambdaworks.FeltFromHex(result_x.Text(16))
	felt_result_y := lambdaworks.FeltFromHex(result_y.Text(16))

	ec.cache[x_addr] = felt_result_x
	ec.cache[x_addr.AddUint(1)] = felt_result_y

	if index-uint(ec.n_input_cells) == 0 {
		return memory.NewMaybeRelocatableFelt(felt_result_x), nil
	} else {
		return memory.NewMaybeRelocatableFelt(felt_result_y), nil
	}
}

func (ec *EcOpBuiltinRunner) ec_op_impl(partial_sum PartialSum, double_point DoublePoint, m lambdaworks.Felt, alpha *big.Int, prime *big.Int, height uint32) (big.Int, big.Int) {

	return *big.NewInt(1), *big.NewInt(1)
}

func (ec *EcOpBuiltinRunner) point_on_curve(x lambdaworks.Felt, y lambdaworks.Felt, alpha lambdaworks.Felt, beta lambdaworks.Felt) bool {
	yp := y.PowUint(2)
	xp := x.PowUint(3).Add(alpha.Mul(x)).Add(beta)
	return yp == xp
}
