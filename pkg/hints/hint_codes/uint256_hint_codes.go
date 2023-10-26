package hint_codes

const UINT256_ADD = "sum_low = ids.a.low + ids.b.low\nids.carry_low = 1 if sum_low >= ids.SHIFT else 0\nsum_high = ids.a.high + ids.b.high + ids.carry_low\nids.carry_high = 1 if sum_high >= ids.SHIFT else 0"
const UINT256_ADD_LOW = "sum_low = ids.a.low + ids.b.low\nids.carry_low = 1 if sum_low >= ids.SHIFT else 0"
const SPLIT_64 = "ids.low = ids.a & ((1<<64) - 1)\nids.high = ids.a >> 64"
const UINT256_SQRT = "from starkware.python.math_utils import isqrt\nn = (ids.n.high << 128) + ids.n.low\nroot = isqrt(n)\nassert 0 <= root < 2 ** 128\nids.root.low = root\nids.root.high = 0"
const UINT256_SQRT_FELT = "from starkware.python.math_utils import isqrt\nn = (ids.n.high << 128) + ids.n.low\nroot = isqrt(n)\nassert 0 <= root < 2 ** 128\nids.root = root;"
const UINT256_SIGNED_NN = "memory[ap] = 1 if 0 <= (ids.a.high % PRIME) < 2 ** 127 else 0"
const UINT256_UNSIGNED_DIV_REM = "a = (ids.a.high << 128) + ids.a.low\ndiv = (ids.div.high << 128) + ids.div.low\nquotient, remainder = divmod(a, div)\n\nids.quotient.low = quotient & ((1 << 128) - 1)\nids.quotient.high = quotient >> 128\nids.remainder.low = remainder & ((1 << 128) - 1)\nids.remainder.high = remainder >> 128"
const UINT256_EXPANDED_UNSIGNED_DIV_REM = "a = (ids.a.high << 128) + ids.a.low\ndiv = (ids.div.b23 << 128) + ids.div.b01\nquotient, remainder = divmod(a, div)\n\nids.quotient.low = quotient & ((1 << 128) - 1)\nids.quotient.high = quotient >> 128\nids.remainder.low = remainder & ((1 << 128) - 1)\nids.remainder.high = remainder >> 128"
const UINT256_MUL_DIV_MOD = "a = (ids.a.high << 128) + ids.a.low\nb = (ids.b.high << 128) + ids.b.low\ndiv = (ids.div.high << 128) + ids.div.low\nquotient, remainder = divmod(a * b, div)\n\nids.quotient_low.low = quotient & ((1 << 128) - 1)\nids.quotient_low.high = (quotient >> 128) & ((1 << 128) - 1)\nids.quotient_high.low = (quotient >> 256) & ((1 << 128) - 1)\nids.quotient_high.high = quotient >> 384\nids.remainder.low = remainder & ((1 << 128) - 1)\nids.remainder.high = remainder >> 128"
const UINT256_SUB = `def split(num: int, num_bits_shift: int = 128, length: int = 2):
    a = []
    for _ in range(length):
        a.append( num & ((1 << num_bits_shift) - 1) )
        num = num >> num_bits_shift
    return tuple(a)

def pack(z, num_bits_shift: int = 128) -> int:
    limbs = (z.low, z.high)
    return sum(limb << (num_bits_shift * i) for i, limb in enumerate(limbs))

a = pack(ids.a)
b = pack(ids.b)
res = (a - b)%2**256
res_split = split(res)
ids.res.low = res_split[0]
ids.res.high = res_split[1]`
