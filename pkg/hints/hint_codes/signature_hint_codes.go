package hint_codes

const DIV_MOD_N_PACKED_DIVMOD_V1 = `from starkware.cairo.common.cairo_secp.secp_utils import N, pack
from starkware.python.math_utils import div_mod, safe_div

a = pack(ids.a, PRIME)
b = pack(ids.b, PRIME)
value = res = div_mod(a, b, N)`

const DIV_MOD_N_PACKED_DIVMOD_EXTERNAL_N = `from starkware.cairo.common.cairo_secp.secp_utils import pack
from starkware.python.math_utils import div_mod, safe_div

a = pack(ids.a, PRIME)
b = pack(ids.b, PRIME)
value = res = div_mod(a, b, N)`

const DIV_MOD_N_SAFE_DIV = "value = k = safe_div(res * b - a, N)"

const DIV_MOD_N_SAFE_DIV_PLUS_ONE = "value = k_plus_one = safe_div(res * b - a, N) + 1"

const XS_SAFE_DIV = "value = k = safe_div(res * s - x, N)"
