package hint_codes

const IMPORT_SECP256R1_ALPHA = "from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_ALPHA as ALPHA"

const IMPORT_SECP256R1_N = "from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_N as N"

const IMPORT_SECP256R1_P = "from starkware.cairo.common.cairo_secp.secp256r1_utils import SECP256R1_P as SECP_P"

const VERIFY_ZERO_EXTERNAL_SECP = "from starkware.cairo.common.cairo_secp.secp_utils import pack\n\nq, r = divmod(pack(ids.val, PRIME), SECP_P)\nassert r == 0, f\"verify_zero: Invalid input {ids.val.d0, ids.val.d1, ids.val.d2}.\"\nids.q = q % PRIME"

const REDUCE_V1 = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack

value = pack(ids.x, PRIME) % SECP_P`

const REDUCE_V2 = `from starkware.cairo.common.cairo_secp.secp_utils import pack
value = pack(ids.x, PRIME) % SECP_P`

const REDUCE_ED25519 = `from starkware.cairo.common.cairo_secp.secp_utils import pack
SECP_P=2**255-19

value = pack(ids.x, PRIME) % SECP_P`

const VERIFY_ZERO_V1 = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack

q, r = divmod(pack(ids.val, PRIME), SECP_P)
assert r == 0, f"verify_zero: Invalid input {ids.val.d0, ids.val.d1, ids.val.d2}."
ids.q = q % PRIME`

const VERIFY_ZERO_V2 = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P
q, r = divmod(pack(ids.val, PRIME), SECP_P)
assert r == 0, f"verify_zero: Invalid input {ids.val.d0, ids.val.d1, ids.val.d2}."
ids.q = q % PRIME`

const VERIFY_ZERO_V3 = `from starkware.cairo.common.cairo_secp.secp_utils import pack
SECP_P = 2**255-19
to_assert = pack(ids.val, PRIME)
q, r = divmod(pack(ids.val, PRIME), SECP_P)
assert r == 0, f"verify_zero: Invalid input {ids.val.d0, ids.val.d1, ids.val.d2}."
ids.q = q % PRIME`

const BIGINT_TO_UINT256 = "ids.low = (ids.x.d0 + ids.x.d1 * ids.BASE) & ((1 << 128) - 1)"

const IS_ZERO_NONDET = "memory[ap] = to_felt_or_relocatable(x == 0)"

const IS_ZERO_INT = "memory[ap] = int(x == 0)"

const IS_ZERO_PACK_V1 = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack

x = pack(ids.x, PRIME) % SECP_P`

const IS_ZERO_PACK_V2 = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack
x = pack(ids.x, PRIME) % SECP_P`
