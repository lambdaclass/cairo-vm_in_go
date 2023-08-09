use starknet_crypto::{poseidon_permute_comp, FieldElement};
extern crate libc;

// C representation of a limbs array: a raw pointer to a mutable unsigned 64 bits integer.
type Limbs = *mut u64;
// C representation of an array of felts: a raw pointer to Limbs.
type PoseidonState = *mut Limbs;

fn FieldElementFromLimbs(limbs: Limbs) -> FieldElement {
    let array = unsafe {
        let slice: &mut [u64] = std::slice::from_raw_parts_mut(limbs, 4);
        let array: [u64; 4] = slice.try_into().unwrap();
        array
    }
    FieldElement::from_mont(array)
}

fn poseidon_permute(state: PoseidonState) {
    let slice = unsafe {
        let slice: &mut [Limbs] = std::slice::from_raw_parts_mut(state, 3);
        slice
    };

}
