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

const GET_POINT_FROM_X = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack

x_cube_int = pack(ids.x_cube, PRIME) % SECP_P
y_square_int = (x_cube_int + ids.BETA) % SECP_P
y = pow(y_square_int, (SECP_P + 1) // 4, SECP_P)

# We need to decide whether to take y or SECP_P - y.
if ids.v % 2 == y % 2:
    value = y
else:
    value = (-y) % SECP_P`
