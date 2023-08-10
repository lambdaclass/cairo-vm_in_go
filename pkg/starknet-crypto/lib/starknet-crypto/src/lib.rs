//use starknet_crypto::{poseidon_permute_comp, FieldElement};
use starknet_crypto::FieldElement;
//extern crate libc;

// C representation of a limbs array: a raw pointer to a mutable unsigned 64 bits integer.
// Consists of 4 u64 values representing a felt in big endian montgomery representation
type Limbs = *mut u64;
// C representation of an array of felts: a raw pointer to Limbs.
type PoseidonState = *mut Limbs;

fn field_element_from_limbs(limbs: Limbs) -> FieldElement {
    let array = unsafe {
        let slice: &mut [u64] = std::slice::from_raw_parts_mut(limbs, 4);
        let array: [u64; 4] = slice.try_into().unwrap();
        array
    };
    FieldElement::from_mont(array)
}

fn limbs_from_field_element(felt: FieldElement, limbs : Limbs) {
    let limb_array = felt.into_mont();
    for i in 0..4 {
        unsafe {
            *limbs.offset(i) = limb_array[i as usize];
        }
    }
}

fn poseidon_permute(state: PoseidonState) {
    // Convert state from C representation to FieldElement
    let slice = unsafe {
        let slice: &mut [Limbs] = std::slice::from_raw_parts_mut(state, 3);
        slice
    };
    let mut state_array =  [FieldElement::ZERO; 3];
    for limbs in slice.iter().take(3_usize) {
        state_array[0] = field_element_from_limbs(*limbs)
    }
    // Call poseidon permute comp
    // Convert state from FieldElement to C representation
    for i in 0..3 {
        unsafe {
            limbs_from_field_element(state_array[i], *state.offset(i as isize))
        }
    }
}
