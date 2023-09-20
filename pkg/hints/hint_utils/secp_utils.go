package hint_utils

import "math/big"

func SECP_P() big.Int {
	secp_p, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)
	return *secp_p
}

func ALPHA() big.Int {
	alpha := big.NewInt(0)
	return *alpha
}

func SECP256R1_ALPHA() big.Int {
	secp_p_alpha, _ := new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853948", 10)
	return *secp_p_alpha
}

func SECP256R1_N() big.Int {
	secp256, _ := new(big.Int).SetString("115792089210356248762697446949407573529996955224135760342422259061068512044369", 10)
	return *secp256
}
