package starknet_crypto_test

import (
	"reflect"
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	starknet_crypto "github.com/lambdaclass/cairo-vm.go/pkg/starknet-crypto"
)

func TestPoseidonPermuteCompA(t *testing.T) {
	// Set initial state values
	poseidon_state := [3]lambdaworks.Felt{
		lambdaworks.FeltFromUint64(3),
		lambdaworks.FeltZero(),
		lambdaworks.FeltFromUint64(2),
	}
	// Run the poseidon permutation
	starknet_crypto.PoseidonPermuteComp(&poseidon_state)
	// Check final state values
	expected_poseidon_state := [3]lambdaworks.Felt{
		lambdaworks.FeltFromHex("0x268c44203f1c763bca21beb5aec78b9063cdcdd0fdf6b598bb8e1e8f2b6253f"),
		lambdaworks.FeltFromHex("0x2b85c9f686f5d3036db55b2ca58a763a3065bc1bc8efbe0e70f3a7171f6cad3"),
		lambdaworks.FeltFromHex("0x61df3789eef0e1ee0dbe010582a00dd099191e6395dfb976e7be3be2fa9d54b"),
	}
	if !reflect.DeepEqual(poseidon_state, expected_poseidon_state) {
		t.Errorf("Wrong state after poseidon permutation.\n Expected %+v.\n Got: %+v", expected_poseidon_state, poseidon_state)
	}
}

func TestPoseidonPermuteCompB(t *testing.T) {
	// Set initial state values
	poseidon_state := [3]lambdaworks.Felt{
		lambdaworks.FeltFromHex("0x268c44203f1c763bca21beb5aec78b9063cdcdd0fdf6b598bb8e1e8f2b6253f"),
		lambdaworks.FeltFromHex("0x2b85c9f686f5d3036db55b2ca58a763a3065bc1bc8efbe0e70f3a7171f6cad3"),
		lambdaworks.FeltFromHex("0x61df3789eef0e1ee0dbe010582a00dd099191e6395dfb976e7be3be2fa9d54b"),
	}
	// Run the poseidon permutation
	starknet_crypto.PoseidonPermuteComp(&poseidon_state)
	// Check final state values
	expected_poseidon_state := [3]lambdaworks.Felt{
		lambdaworks.FeltFromHex("0x4ec565b1b01606b5222602b20f8ddc4a8a7c75b559b852ab183a0daf5930b5c"),
		lambdaworks.FeltFromHex("0x4d3c32c3c7cd39b6444db42e2437eeda12e459d28ce49a0f761a23d64c29e4c"),
		lambdaworks.FeltFromHex("0x749d4d0ddf41548e039f183b745a08b80fad54e9ac389021148350bdda70a92"),
	}
	if !reflect.DeepEqual(poseidon_state, expected_poseidon_state) {
		t.Errorf("Wrong state after poseidon permutation.\n Expected %+v.\n Got: %+v", expected_poseidon_state, poseidon_state)
	}
}
