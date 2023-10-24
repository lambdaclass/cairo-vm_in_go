package hint_codes

const BLAKE2S_COMPUTE = `from starkware.cairo.common.cairo_blake2s.blake2s_utils import compute_blake2s_func
compute_blake2s_func(segments=segments, output_ptr=ids.output)`
