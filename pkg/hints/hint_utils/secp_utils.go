package hint_utils

import (
	"errors"
	"math/big"
)

func SECP_P() big.Int {
	secpP, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)
	return *secpP
}

func SECP_P_V2() big.Int {
	secpP, _ := new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)
	return *secpP
}

func ALPHA() big.Int {
	alpha := big.NewInt(0)
	return *alpha
}

func SECP256R1_ALPHA() big.Int {
	secpPalpha, _ := new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853948", 10)
	return *secpPalpha
}

func SECP256R1_N() big.Int {
	secp256, _ := new(big.Int).SetString("115792089210356248762697446949407573529996955224135760342422259061068512044369", 10)
	return *secp256
}

func SECP256R1_P() big.Int {
	secp256r1, _ := new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	return *secp256r1
}

func BASE_MINUS_ONE() *big.Int {
	res, _ := new(big.Int).SetString("77371252455336267181195263", 10)
	return res
}

func Bigint3Split(integer big.Int) ([]big.Int, error) {
	canonicalRepr := make([]big.Int, 3)
	num := integer

	for i := 0; i < 3; i++ {
		canonicalRepr[i] = *new(big.Int).And(&num, BASE_MINUS_ONE())
		num = *new(big.Int).Rsh(&num, 86)
	}
	if num.Cmp(big.NewInt(0)) != 0 {
		return nil, errors.New("HintError SecpSplitOutOfRange")
	}

	return canonicalRepr, nil
}
