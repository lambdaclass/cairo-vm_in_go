package hint_codes

const EC_NEGATE = "from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack\n\ny = pack(ids.point.y, PRIME) % SECP_P\n# The modulo operation in python always returns a nonnegative number.\nvalue = (-y) % SECP_P"
const EC_NEGATE_EMBEDDED_SECP = "from starkware.cairo.common.cairo_secp.secp_utils import pack\nSECP_P = 2**255-19\n\ny = pack(ids.point.y, PRIME) % SECP_P\n# The modulo operation in python always returns a nonnegative number.\nvalue = (-y) % SECP_P"
const EC_MUL_INNER = "memory[ap] = (ids.scalar % PRIME) % 2"
