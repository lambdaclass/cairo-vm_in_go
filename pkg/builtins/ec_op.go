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

type PartialSumB struct {
	x big.Int
	y big.Int
}

type DoublePointB struct {
	x big.Int
	y big.Int
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
	EC_POINT_INDICES := [3]EcPoint{{x: 0, y: 1}, {x: 2, y: 3}, {x: 5, y: 6}}
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
		if !PointOnCurve(x, y, alpha, beta) {
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

	result, err := ec.EcOnImpl(partial_sum, double_point, input_cells[4], alpha_big_int, prime, ec.scalar_height)

	felt_result_x := lambdaworks.FeltFromHex(result.x.Text(16))
	felt_result_y := lambdaworks.FeltFromHex(result.y.Text(16))

	ec.cache[x_addr] = felt_result_x
	ec.cache[x_addr.AddUint(1)] = felt_result_y

	if index-uint(ec.n_input_cells) == 0 {
		return memory.NewMaybeRelocatableFelt(felt_result_x), nil
	} else {
		return memory.NewMaybeRelocatableFelt(felt_result_y), nil
	}
}

func LineSlope(point_a PartialSumB, point_b DoublePointB, prime big.Int) (big.Int, error) {
	mod_value := new(big.Int).Sub(&point_a.x, &point_b.y)
	mod_value.Mod(mod_value, &prime)

	if mod_value == big.NewInt(0) {
		return big.Int{}, errors.New("is multiple of prime")
	}

	n := new(big.Int).Sub(&point_a.y, &point_b.y)
	m := new(big.Int).Sub(&point_a.x, &point_b.x)

	z, _ := new(big.Int).DivMod(n, m, &prime)

	return *z, nil
}

func EcAdd(point_a PartialSumB, point_b DoublePointB, prime big.Int) (PartialSumB, error) {
	m, err := LineSlope(point_a, point_b, prime)
	if err != nil {
		return PartialSumB{}, err
	}

	x := new(big.Int).Mul(&m, &m)
	x.Sub(&point_a.x, &point_b.x)
	x.Mod(x, &prime)

	y := new(big.Int).Mul(&m, new(big.Int).Sub(&point_a.x, x))
	y.Sub(y, &point_a.y)
	y.Mod(y, &prime)

	return PartialSumB{x: *x, y: *y}, nil
}

func EcDoubleSlope(point DoublePointB, alpha big.Int, prime big.Int) (big.Int, error) {
	q := new(big.Int).Mod(&point.y, &prime)
	if q == big.NewInt(0) {
		return big.Int{}, errors.New("is multiple of prime")
	}

	n := new(big.Int).Mul(&point.x, &point.x)
	n.Mul(n, big.NewInt(3))
	n.Add(n, &alpha)

	m := new(big.Int).Mul(&point.y, big.NewInt(2))

	z, _ := new(big.Int).DivMod(n, m, &prime)

	return *z, nil
}

func ec_double(point DoublePointB, alpha big.Int, prime big.Int) (DoublePointB, error) {
	m, err := EcDoubleSlope(point, alpha, prime)
	if err != nil {
		return DoublePointB{}, err
	}

	x := new(big.Int).Mul(&m, &m)
	x.Sub(x, new(big.Int).Mul(big.NewInt(2), &point.x))
	x.Mod(x, &prime)

	y := new(big.Int).Mul(&m, new(big.Int).Sub(&point.x, x))
	y.Sub(y, &point.y)
	y.Mod(y, &prime)

	return DoublePointB{x: *x, y: *y}, nil
}

func (ec *EcOpBuiltinRunner) EcOnImpl(partial_sum PartialSum, double_point DoublePoint, m lambdaworks.Felt, alpha *big.Int, prime *big.Int, height uint32) (PartialSumB, error) {
	slope, _ := m.ToBigInt()
	partial_sum_b_x, _ := partial_sum.x.ToBigInt()
	partial_sum_b_y, _ := partial_sum.y.ToBigInt()

	partial_sum_b := PartialSumB{x: partial_sum_b_x, y: partial_sum_b_y}

	double_point_b_x, _ := double_point.x.ToBigInt()
	double_point_b_y, _ := double_point.y.ToBigInt()

	double_point_b := DoublePointB{x: double_point_b_x, y: double_point_b_y}

	for i := 0; i < int(height); i++ {
		var err error
		if (double_point_b.x.Sub(&double_point_b.x, &partial_sum_b.x)) == big.NewInt(0) {
			return PartialSumB{}, errors.New("Runner error EcOpSameXCoordinate")
		}
		if !((slope.And(&slope, big.NewInt(1))) == big.NewInt(0)) {
			partial_sum_b, err = EcAdd(partial_sum_b, double_point_b, *prime)
			if err != nil {
				return PartialSumB{}, err
			}
		}
		double_point_b, err = ec_double(double_point_b, *alpha, *prime)
		slope = *slope.Rsh(&slope, 1)
	}

	return partial_sum_b, nil
}

func PointOnCurve(x lambdaworks.Felt, y lambdaworks.Felt, alpha lambdaworks.Felt, beta lambdaworks.Felt) bool {
	yp := y.PowUint(2)
	xp := x.PowUint(3).Add(alpha.Mul(x)).Add(beta)
	return yp == xp
}
