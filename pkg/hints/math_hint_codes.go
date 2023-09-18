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

const ASSERT_NOT_EQUAL = "from starkware.cairo.lang.vm.relocatable import RelocatableValue\nboth_ints = isinstance(ids.a, int) and isinstance(ids.b, int)\nboth_relocatable = (\n    isinstance(ids.a, RelocatableValue) and isinstance(ids.b, RelocatableValue) and\n    ids.a.segment_index == ids.b.segment_index)\nassert both_ints or both_relocatable, \\\n    f'assert_not_equal failed: non-comparable values: {ids.a}, {ids.b}.'\nassert (ids.a - ids.b) % PRIME != 0, f'assert_not_equal failed: {ids.a} = {ids.b}.'"

const SQRT = "from starkware.python.math_utils import isqrt\nvalue = ids.value % PRIME\nassert value < 2 ** 250, f\"value={value} is outside of the range [0, 2**250).\"\nassert 2 ** 250 < PRIME\nids.root = isqrt(value)"

const ASSERT_LE_FELT = "import itertools\n\nfrom starkware.cairo.common.math_utils import assert_integer\nassert_integer(ids.a)\nassert_integer(ids.b)\na = ids.a % PRIME\nb = ids.b % PRIME\nassert a <= b, f'a = {a} is not less than or equal to b = {b}.'\n\n# Find an arc less than PRIME / 3, and another less than PRIME / 2.\nlengths_and_indices = [(a, 0), (b - a, 1), (PRIME - 1 - b, 2)]\nlengths_and_indices.sort()\nassert lengths_and_indices[0][0] <= PRIME // 3 and lengths_and_indices[1][0] <= PRIME // 2\nexcluded = lengths_and_indices[2][1]\n\nmemory[ids.range_check_ptr + 1], memory[ids.range_check_ptr + 0] = (\n    divmod(lengths_and_indices[0][0], ids.PRIME_OVER_3_HIGH))\nmemory[ids.range_check_ptr + 3], memory[ids.range_check_ptr + 2] = (\n    divmod(lengths_and_indices[1][0], ids.PRIME_OVER_2_HIGH))"
