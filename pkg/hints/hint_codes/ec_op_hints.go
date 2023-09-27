package hint_codes

const EC_NEGATE = "from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack\n\ny = pack(ids.point.y, PRIME) % SECP_P\n# The modulo operation in python always returns a nonnegative number.\nvalue = (-y) % SECP_P"
const EC_NEGATE_EMBEDDED_SECP = "from starkware.cairo.common.cairo_secp.secp_utils import pack\nSECP_P = 2**255-19\n\ny = pack(ids.point.y, PRIME) % SECP_P\n# The modulo operation in python always returns a nonnegative number.\nvalue = (-y) % SECP_P"
const EC_DOUBLE_SLOPE_V1 = "from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack\nfrom starkware.python.math_utils import ec_double_slope\n\n# Compute the slope.\nx = pack(ids.point.x, PRIME)\ny = pack(ids.point.y, PRIME)\nvalue = slope = ec_double_slope(point=(x, y), alpha=0, p=SECP_P)"
const COMPUTE_SLOPE_V1 = "from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack\nfrom starkware.python.math_utils import line_slope\n\n# Compute the slope.\nx0 = pack(ids.point0.x, PRIME)\ny0 = pack(ids.point0.y, PRIME)\nx1 = pack(ids.point1.x, PRIME)\ny1 = pack(ids.point1.y, PRIME)\nvalue = slope = line_slope(point1=(x0, y0), point2=(x1, y1), p=SECP_P)"
const EC_DOUBLE_ASSIGN_NEW_X_V1 = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack

slope = pack(ids.slope, PRIME)
x = pack(ids.point.x, PRIME)
y = pack(ids.point.y, PRIME)

value = new_x = (pow(slope, 2, SECP_P) - 2 * x) % SECP_P`
const EC_DOUBLE_ASSIGN_NEW_X_V2 = `from starkware.cairo.common.cairo_secp.secp_utils import pack

slope = pack(ids.slope, PRIME)
x = pack(ids.point.x, PRIME)
y = pack(ids.point.y, PRIME)

value = new_x = (pow(slope, 2, SECP_P) - 2 * x) % SECP_P`
const EC_DOUBLE_ASSIGN_NEW_X_V3 = `from starkware.cairo.common.cairo_secp.secp_utils import pack
SECP_P = 2**255-19

slope = pack(ids.slope, PRIME)
x = pack(ids.point.x, PRIME)
y = pack(ids.point.y, PRIME)

value = new_x = (pow(slope, 2, SECP_P) - 2 * x) % SECP_P`
const EC_DOUBLE_ASSIGN_NEW_X_V4 = `from starkware.cairo.common.cairo_secp.secp_utils import SECP_P, pack

slope = pack(ids.slope, PRIME)
x = pack(ids.pt.x, PRIME)
y = pack(ids.pt.y, PRIME)

value = new_x = (pow(slope, 2, SECP_P) - 2 * x) % SECP_P`
const EC_DOUBLE_SLOPE_EXTERNAL_CONSTS = "from starkware.cairo.common.cairo_secp.secp_utils import pack\nfrom starkware.python.math_utils import ec_double_slope\n\n# Compute the slope.\nx = pack(ids.point.x, PRIME)\ny = pack(ids.point.y, PRIME)\nvalue = slope = ec_double_slope(point=(x, y), alpha=ALPHA, p=SECP_P)"
const NONDET_BIGINT3_V1 = "from starkware.cairo.common.cairo_secp.secp_utils import split\n\nsegments.write_arg(ids.res.address_, split(value))"
