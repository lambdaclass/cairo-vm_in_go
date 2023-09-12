package hints

const ASSERT_NN = "from starkware.cairo.common.math_utils import assert_integer\nassert_integer(ids.a)\nassert 0 <= ids.a % PRIME < range_check_builtin.bound, f'a = {ids.a} is out of range.'"

const VERIFY_ECDSA_SIGNATURE = "ecdsa_builtin.add_signature(ids.ecdsa_ptr.address_, (ids.signature_r, ids.signature_s))"
