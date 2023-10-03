package hint_codes

const BLAKE2S_COMPUTE = `from starkware.cairo.common.cairo_blake2s.blake2s_utils import compute_blake2s_func
compute_blake2s_func(segments=segments, output_ptr=ids.output)`

const BLAKE2S_ADD_UINT256 = `B = 32
MASK = 2 ** 32 - 1
segments.write_arg(ids.data, [(ids.low >> (B * i)) & MASK for i in range(4)])
segments.write_arg(ids.data + 4, [(ids.high >> (B * i)) & MASK for i in range(4)])`

const BLAKE2S_ADD_UINT256_BIGEND = `B = 32
MASK = 2 ** 32 - 1
segments.write_arg(ids.data, [(ids.high >> (B * (3 - i))) & MASK for i in range(4)])
segments.write_arg(ids.data + 4, [(ids.low >> (B * (3 - i))) & MASK for i in range(4)])`

const BLAKE2S_FINALIZE = `# Add dummy pairs of input and output.
from starkware.cairo.common.cairo_blake2s.blake2s_utils import IV, blake2s_compress

_n_packed_instances = int(ids.N_PACKED_INSTANCES)
assert 0 <= _n_packed_instances < 20
_blake2s_input_chunk_size_felts = int(ids.INPUT_BLOCK_FELTS)
assert 0 <= _blake2s_input_chunk_size_felts < 100

message = [0] * _blake2s_input_chunk_size_felts
modified_iv = [IV[0] ^ 0x01010020] + IV[1:]
output = blake2s_compress(
    message=message,
    h=modified_iv,
    t0=0,
    t1=0,
    f0=0xffffffff,
    f1=0,
)
padding = (modified_iv + message + [0, 0xffffffff] + output) * (_n_packed_instances - 1)
segments.write_arg(ids.blake2s_ptr_end, padding)`
