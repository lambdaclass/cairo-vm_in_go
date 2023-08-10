use starknet_crypto::{poseidon_permute_comp, FieldElement};
extern crate libc;

// C representation of a limbs array: a raw pointer to a mutable unsigned 64 bits integer.
// Consists of 4 u64 values representing a felt in big endian montgomery representation
type Limbs = *mut u64;

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

#[no_mangle]
extern "C" fn poseidon_permute(first_state_felt: Limbs, second_state_felt: Limbs, third_state_felt: Limbs) {
    // Convert state from C representation to FieldElement
    let mut state_array =  [FieldElement::ZERO; 3];
    state_array[0] = field_element_from_limbs(first_state_felt);
    state_array[1] = field_element_from_limbs(second_state_felt);
    state_array[2] = field_element_from_limbs(third_state_felt);
    // Call poseidon permute comp
    poseidon_permute_comp(&mut state_array);
    // Convert state from FieldElement back to C representation
    limbs_from_field_element(state_array[0], first_state_felt);
    limbs_from_field_element(state_array[1], second_state_felt);
    limbs_from_field_element(state_array[2], third_state_felt);
}
