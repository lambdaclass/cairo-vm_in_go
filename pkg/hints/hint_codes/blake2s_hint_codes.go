package hint_codes

const BLAKE2S_COMPUTE = `from starkware.cairo.common.cairo_blake2s.blake2s_utils import compute_blake2s_func
compute_blake2s_func(segments=segments, output_ptr=ids.output)`

const BLAKE2S_ADD_UINT256_BIGEND = `B = 32
MASK = 2 ** 32 - 1
segments.write_arg(ids.data, [(ids.high >> (B * (3 - i))) & MASK for i in range(4)])
segments.write_arg(ids.data + 4, [(ids.low >> (B * (3 - i))) & MASK for i in range(4)])`
