package hint_utils

import (
	"math/big"

	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

type Uint256 struct {
	Low  Felt
	High Felt
}

func (ui256 *Uint256) ToString() string {
	return "Uint256 {low: " + ui256.Low.ToSignedFeltString() + ", high: " + ui256.High.ToSignedFeltString() + "}"
}

/*
Returns a Uint256 as a big.Int

	res = high << 128 + low
*/
func (u *Uint256) ToBigInt() *big.Int {
	high := new(big.Int).Lsh(u.High.ToBigInt(), 128)
	low := u.Low.ToBigInt()
	res := new(big.Int).Add(high, low)
	return res
}

/*
Returns a big.Int as Uint256
*/
func ToUint256(a *big.Int) Uint256 {
	maxU128, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
	low := new(big.Int).And(a, maxU128)
	high := new(big.Int).Rsh(a, 128)
	return Uint256{Low: FeltFromBigInt(low), High: FeltFromBigInt(high)}
}

func (u *Uint256) IsEqual(other Uint256) bool {
	return u.Low.Cmp(other.Low) == 0 && u.High.Cmp(other.High) == 0
}
