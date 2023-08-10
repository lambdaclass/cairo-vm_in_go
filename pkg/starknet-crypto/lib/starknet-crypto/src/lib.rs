use starknet_crypto::{poseidon_permute_comp, FieldElement};
extern crate libc;

// C representation of a bit array: a raw pointer to a mutable unsigned 8 bits integer.
type Bytes = *mut u8;

fn field_element_from_bytes(bytes: Bytes) -> FieldElement {
    let array = unsafe {
        let slice: &mut [u8] = std::slice::from_raw_parts_mut(bytes, 32);
        let array: [u8; 32] = slice.try_into().unwrap();
        array
    };
    FieldElement::from_bytes_be(&array).unwrap()
}

fn bytes_from_field_element(felt: FieldElement, bytes : Bytes) {
    let byte_array = felt.to_bytes_be();
    for i in 0..32 {
        unsafe {
            *bytes.offset(i) = byte_array[i as usize];
        }
    }
}

#[no_mangle]
extern "C" fn poseidon_permute(first_state_felt: Bytes, second_state_felt: Bytes, third_state_felt: Bytes) {
    // Convert state from C representation to FieldElement
    let mut state_array =  [FieldElement::ZERO; 3];
    state_array[0] = field_element_from_bytes(first_state_felt);
    state_array[1] = field_element_from_bytes(second_state_felt);
    state_array[2] = field_element_from_bytes(third_state_felt);
    println!("State array {:?}, {:?}, {:?}", state_array[0], state_array[1], state_array[2]);
    // Call poseidon permute comp
    poseidon_permute_comp(&mut state_array);
    // Convert state from FieldElement back to C representation
    bytes_from_field_element(state_array[0], first_state_felt);
    bytes_from_field_element(state_array[1], second_state_felt);
    bytes_from_field_element(state_array[2], third_state_felt);
}
