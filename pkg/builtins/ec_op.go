package builtins

import (
	"errors"
	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/math_utils"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
	"math/big"
	"reflect"
)

type EcOpBuiltinRunner struct {
	included      bool
	base          memory.Relocatable
	cache         map[memory.Relocatable]lambdaworks.Felt
	scalar_height uint32
}

type EcPoint struct {
	x uint
	y uint
}

type PartialSum struct {
	X lambdaworks.Felt
	Y lambdaworks.Felt
}

type DoublePoint struct {
	X lambdaworks.Felt
	Y lambdaworks.Felt
}

type PartialSumB struct {
	X big.Int
	Y big.Int
}

type DoublePointB struct {
	X big.Int
	Y big.Int
}

const INPUT_CELLS_PER_EC_OP = 5
const CELLS_PER_EC_OP = 7
const EC_OP_BUILTIN_NAME = "ec_op"

const PRIME = "0x800000000000011000000000000000000000000000000000000000000000001"

func NewEcOpBuiltinRunner() *EcOpBuiltinRunner {
	return &EcOpBuiltinRunner{
		cache:         make(map[memory.Relocatable]lambdaworks.Felt),
		scalar_height: 256,
	}
}

func (ec *EcOpBuiltinRunner) AddValidationRule(*memory.Memory) {}

func (ec *EcOpBuiltinRunner) Base() memory.Relocatable {
	return ec.base
}

func (ec *EcOpBuiltinRunner) Name() string {
	return EC_OP_BUILTIN_NAME
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

func (ec *EcOpBuiltinRunner) Include(include bool) {
	ec.included = include
}

func (ec *EcOpBuiltinRunner) DeduceMemoryCell(address memory.Relocatable, mem *memory.Memory) (*memory.MaybeRelocatable, error) {
	EC_POINT_INDICES := [3]EcPoint{{x: 0, y: 1}, {x: 2, y: 3}, {x: 5, y: 6}}
	OUTPUT_INDICES := EC_POINT_INDICES[2]
	alpha := lambdaworks.FeltOne()
	beta_low := lambdaworks.FeltFromHex("0x609ad26c15c915c1f4cdfcb99cee9e89")
	beta_high := lambdaworks.FeltFromHex("0x6f21413efbe40de150e596d72f7a8c5")
	beta := (beta_high.Shl(128)).Add(beta_low)

	index := address.Offset % uint(CELLS_PER_EC_OP)

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

	input_cells := make([]lambdaworks.Felt, 0)
	for i := 0; i < int(INPUT_CELLS_PER_EC_OP); i++ {
		maybe_rel, err := mem.Get(instance.AddUint(uint(i)))
		if err == nil {
			felt, is_felt := maybe_rel.GetFelt()
			if is_felt {
				input_cells = append(input_cells, felt)
			} else {
				return nil, errors.New("Runner error, Expected Integer for input cells")
			}
		} else {
			return nil, nil
		}
	}

	for j := 0; j < 2; j++ {
		x := input_cells[EC_POINT_INDICES[j].x]
		y := input_cells[EC_POINT_INDICES[j].y]
		if !PointOnCurve(x, y, alpha, beta) {
			return nil, errors.New("Point not in curve")
		}
	}

	prime, ok := new(big.Int).SetString(PRIME[2:], 16)
	if !ok {
		return nil, errors.New("Could not parse prime")
	}

	alpha_big_int := big.NewInt(1)

	partial_sum := PartialSum{X: input_cells[0], Y: input_cells[1]}
	double_point := DoublePoint{X: input_cells[2], Y: input_cells[3]}

	result, err := EcOnImpl(partial_sum, double_point, input_cells[4], alpha_big_int, prime, ec.scalar_height)

	felt_result_x := lambdaworks.FeltFromBeBytes((*[32]byte)(result.X.Bytes()))
	felt_result_y := lambdaworks.FeltFromBeBytes((*[32]byte)(result.Y.Bytes()))

	ec.cache[x_addr] = felt_result_x
	ec.cache[x_addr.AddUint(1)] = felt_result_y

	if index-uint(INPUT_CELLS_PER_EC_OP) == 0 {
		return memory.NewMaybeRelocatableFelt(felt_result_x), nil
	} else {
		return memory.NewMaybeRelocatableFelt(felt_result_y), nil
	}
}

func LineSlope(point_a PartialSumB, point_b DoublePointB, prime big.Int) (big.Int, error) {
	mod_value := new(big.Int).Sub(&point_a.X, &point_b.X)
	mod_value.Mod(mod_value, &prime)

	if mod_value.Cmp(big.NewInt(0)) == 0 {
		return big.Int{}, errors.New("is multiple of prime")
	}

	n := new(big.Int).Sub(&point_a.Y, &point_b.Y)
	m := new(big.Int).Sub(&point_a.X, &point_b.X)

	z, _ := math_utils.DivMod(n, m, &prime)
	return *z, nil
}

func EcAdd(point_a PartialSumB, point_b DoublePointB, prime big.Int) (PartialSumB, error) {
	m, err := LineSlope(point_a, point_b, prime)
	if err != nil {
		return PartialSumB{}, err
	}

	x := new(big.Int).Mul(&m, &m)
	x.Sub(x, &point_a.X)
	x.Sub(x, &point_b.X)
	x.Mod(x, &prime)

	y := new(big.Int).Mul(&m, new(big.Int).Sub(&point_a.X, x))
	y.Sub(y, &point_a.Y)
	y.Mod(y, &prime)

	return PartialSumB{X: *x, Y: *y}, nil
}

func EcDoubleSlope(point DoublePointB, alpha big.Int, prime big.Int) (big.Int, error) {
	q := new(big.Int).Mod(&point.Y, &prime)
	if q == big.NewInt(0) {
		return big.Int{}, errors.New("is multiple of prime")
	}
	n := new(big.Int).Mul(&point.X, &point.X)
	n.Mul(n, big.NewInt(3))
	n.Add(n, &alpha)

	m := new(big.Int).Mul(&point.Y, big.NewInt(2))
	z, _ := math_utils.DivMod(n, m, &prime)

	return *z, nil
}

func EcDouble(point DoublePointB, alpha big.Int, prime big.Int) (DoublePointB, error) {
	m, err := EcDoubleSlope(point, alpha, prime)
	if err != nil {
		return DoublePointB{}, err
	}
	x := new(big.Int).Mul(&m, &m)
	x.Sub(x, new(big.Int).Mul(big.NewInt(2), &point.X))
	x.Mod(x, &prime)

	y := new(big.Int).Mul(&m, new(big.Int).Sub(&point.X, x))
	y.Sub(y, &point.Y)
	y.Mod(y, &prime)

	return DoublePointB{X: *x, Y: *y}, nil
}

func EcOnImpl(partial_sum PartialSum, double_point DoublePoint, m lambdaworks.Felt, alpha *big.Int, prime *big.Int, height uint32) (PartialSumB, error) {
	slope := m.ToBigInt()
	partial_sum_b_x := partial_sum.X.ToBigInt()
	partial_sum_b_y := partial_sum.Y.ToBigInt()
	partial_sum_b := PartialSumB{X: *partial_sum_b_x, Y: *partial_sum_b_y}

	double_point_b_x := double_point.X.ToBigInt()
	double_point_b_y := double_point.Y.ToBigInt()
	double_point_b := DoublePointB{X: *double_point_b_x, Y: *double_point_b_y}

	for i := 0; i < int(height); i++ {
		var err error
		if (new(big.Int).Sub(&double_point_b.X, &partial_sum_b.X)).Cmp(big.NewInt(0)) == 0 {
			return PartialSumB{}, errors.New("Runner error EcOpSameXCoordinate")
		}

		and_operation := (new(big.Int).And(slope, big.NewInt(1)))
		if and_operation.Cmp(big.NewInt(0)) > 0 {
			partial_sum_b, err = EcAdd(partial_sum_b, double_point_b, *prime)
			if err != nil {
				return PartialSumB{}, err
			}
		}

		double_point_b, err = EcDouble(double_point_b, *alpha, *prime)
		slope = slope.Rsh(slope, 1)
	}
	return partial_sum_b, nil
}

func PointOnCurve(x lambdaworks.Felt, y lambdaworks.Felt, alpha lambdaworks.Felt, beta lambdaworks.Felt) bool {
	yp := y.PowUint(2)
	xp := x.PowUint(3).Add(alpha.Mul(x)).Add(beta)
	return yp == xp
}
