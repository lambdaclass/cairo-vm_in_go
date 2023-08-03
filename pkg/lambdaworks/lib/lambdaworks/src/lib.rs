use lambdaworks_math::{
    field::element::FieldElement,
    field::fields::fft_friendly::stark_252_prime_field::Stark252PrimeField,
    unsigned_integer::element::UnsignedInteger,
    unsigned_integer::element::U256
};

extern crate libc;
// A 256 bit prime field represented as a Montgomery, 4-limb UnsignedInteger.
type Felt = FieldElement<Stark252PrimeField>;

// C representation of a limbs array: a raw pointer to a mutable unsigned 64 bits integer.
type Limbs = *mut u64;

// Receives a Felt and writes its C representation in the limbs variable, as we can't
// return arrays in C.
//
// Felt uses the montgomery representation internally, so to be able to reconstruct a felt
// in a different call, the representative limbs are the ones written as a result of this call.
fn felt_to_limbs(felt: Felt, limbs: Limbs) {
    let representative = felt.representative().limbs;
    for i in 0..4 {
        let u = i as usize;
        unsafe {
            *limbs.offset(i) = representative[u];
        }
    }
}

// Receives a C representation of a limbs array and returns a felt representing
// the same number.
fn limbs_to_felt(limbs: Limbs) -> Felt {
    unsafe {
        let slice: &mut [u64] = std::slice::from_raw_parts_mut(limbs, 4);
        let array: [u64; 4] = slice.try_into().unwrap();
        let ui = UnsignedInteger::from_limbs(array);
        let felt = Felt::from(&ui);
        return felt;
    }
}

#[no_mangle]
pub extern "C" fn from(result: Limbs, value: u64) {
    felt_to_limbs(Felt::from(value), result);
}

#[no_mangle]
pub extern "C" fn from_hex(result: Limbs, value: *const libc::c_char) {
    let val_cstr = unsafe { core::ffi::CStr::from_ptr(value) };
    let value = val_cstr.to_str().unwrap();
    let felt = match FieldElement::from_hex(value) {
        Ok(felt) => felt,
        Err(_) => {
            panic!("Failed to convert hexadecimal string to FieldElement.");
        }
    };
    felt_to_limbs(felt, result);
}

#[no_mangle]
pub extern "C" fn from_dec_str(result: Limbs, value: *const libc::c_char) {
    let val_cstr = unsafe { core::ffi::CStr::from_ptr(value) };
    let val_str = val_cstr.to_str().unwrap();
    let felt = match val_str.strip_prefix("-") {
        Some(stripped) => {
            let val = U256::from_dec_str(stripped).unwrap(); 
            Felt::from(0) - Felt::from(&val)
        }
        None => { 
            Felt::from(&U256::from_dec_str(val_str).unwrap())
        }
    };
    felt_to_limbs(felt, result)
    
}

#[no_mangle]
pub extern "C" fn zero(result: Limbs) {
    felt_to_limbs(Felt::zero(), result)
}

#[no_mangle]
pub extern "C" fn one(result: Limbs) {
    felt_to_limbs(Felt::one(), result)
}

#[no_mangle]
pub extern "C" fn add(a: Limbs, b: Limbs, result: Limbs) {
    felt_to_limbs(limbs_to_felt(a) + limbs_to_felt(b), result);
}

#[no_mangle]
pub extern "C" fn sub(a: Limbs, b: Limbs, result: Limbs) {
    felt_to_limbs(limbs_to_felt(a) - limbs_to_felt(b), result)
}

#[no_mangle]
pub extern "C" fn mul(a: Limbs, b: Limbs, result: Limbs) {
    felt_to_limbs(limbs_to_felt(a) * limbs_to_felt(b), result)
}

#[no_mangle]
pub extern "C" fn lw_div(a: Limbs, b: Limbs, result: Limbs) {
    felt_to_limbs(limbs_to_felt(a) / limbs_to_felt(b), result)
}
