package hint_codes

const UINT256_ADD = "sum_low = ids.a.low + ids.b.low\nids.carry_low = 1 if sum_low >= ids.SHIFT else 0\nsum_high = ids.a.high + ids.b.high + ids.carry_low\nids.carry_high = 1 if sum_high >= ids.SHIFT else 0"
const UINT256_ADD_LOW = "sum_low = ids.a.low + ids.b.low\nids.carry_low = 1 if sum_low >= ids.SHIFT else 0"
const SPLIT_64 = "ids.low = ids.a & ((1<<64) - 1)\nids.high = ids.a >> 64"
const UINT256_SQRT = "from starkware.python.math_utils import isqrt\nn = (ids.n.high << 128) + ids.n.low\nroot = isqrt(n)\nassert 0 <= root < 2 ** 128\nids.root.low = root\nids.root.high = 0"
const UINT256_SQRT_FELT = "from starkware.python.math_utils import isqrt\nn = (ids.n.high << 128) + ids.n.low\nroot = isqrt(n)\nassert 0 <= root < 2 ** 128\nids.root = root;"
