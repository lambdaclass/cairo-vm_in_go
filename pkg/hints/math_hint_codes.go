package hints

const ASSERT_NN = "from starkware.cairo.common.math_utils import assert_integer\nassert_integer(ids.a)\nassert 0 <= ids.a % PRIME < range_check_builtin.bound, f'a = {ids.a} is out of range.'"

const IS_POSITIVE = "from starkware.cairo.common.math_utils import is_positive\nids.is_positive = 1 if is_positive(\n    value=ids.value, prime=PRIME, rc_bound=range_check_builtin.bound) else 0"

const ASSERT_NOT_ZERO = "from starkware.cairo.common.math_utils import assert_integer\nassert_integer(ids.value)\nassert ids.value % PRIME != 0, f'assert_not_zero failed: {ids.value} = 0.'"

const IS_QUAD_RESIDUE = `from starkware.crypto.signature.signature import FIELD_PRIME
from starkware.python.math_utils import div_mod, is_quad_residue, sqrt

x = ids.x
if is_quad_residue(x, FIELD_PRIME):
    ids.y = sqrt(x, FIELD_PRIME)
else:
    ids.y = sqrt(div_mod(x, 3, FIELD_PRIME), FIELD_PRIME)`
