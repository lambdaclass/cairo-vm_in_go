package hint_codes

const UNSAFE_KECCAK = "from eth_hash.auto import keccak\n\ndata, length = ids.data, ids.length\n\nif '__keccak_max_size' in globals():\n    assert length <= __keccak_max_size, \\\n        f'unsafe_keccak() can only be used with length<={__keccak_max_size}. ' \\\n        f'Got: length={length}.'\n\nkeccak_input = bytearray()\nfor word_i, byte_i in enumerate(range(0, length, 16)):\n    word = memory[data + word_i]\n    n_bytes = min(16, length - byte_i)\n    assert 0 <= word < 2 ** (8 * n_bytes)\n    keccak_input += word.to_bytes(n_bytes, 'big')\n\nhashed = keccak(keccak_input)\nids.high = int.from_bytes(hashed[:16], 'big')\nids.low = int.from_bytes(hashed[16:32], 'big')"

const UNSAFE_KECCAK_FINALIZE = "from eth_hash.auto import keccak\nkeccak_input = bytearray()\nn_elms = ids.keccak_state.end_ptr - ids.keccak_state.start_ptr\nfor word in memory.get_range(ids.keccak_state.start_ptr, n_elms):\n    keccak_input += word.to_bytes(16, 'big')\nhashed = keccak(keccak_input)\nids.high = int.from_bytes(hashed[:16], 'big')\nids.low = int.from_bytes(hashed[16:32], 'big')"

const COMPARE_BYTES_IN_WORD_NONDET = "memory[ap] = to_felt_or_relocatable(ids.n_bytes < ids.BYTES_IN_WORD)"

const COMPARE_KECCAK_FULL_RATE_IN_BYTES_NONDET = "memory[ap] = to_felt_or_relocatable(ids.n_bytes >= ids.KECCAK_FULL_RATE_IN_BYTES)"

const BLOCK_PERMUTATION = `from starkware.cairo.common.keccak_utils.keccak_utils import keccak_func
_keccak_state_size_felts = int(ids.KECCAK_STATE_SIZE_FELTS)
assert 0 <= _keccak_state_size_felts < 100

output_values = keccak_func(memory.get_range(
    ids.keccak_ptr - _keccak_state_size_felts, _keccak_state_size_felts))
segments.write_arg(ids.keccak_ptr, output_values)`

const CAIRO_KECCAK_FINALIZE_V1 = `# Add dummy pairs of input and output.
_keccak_state_size_felts = int(ids.KECCAK_STATE_SIZE_FELTS)
_block_size = int(ids.BLOCK_SIZE)
assert 0 <= _keccak_state_size_felts < 100
assert 0 <= _block_size < 10
inp = [0] * _keccak_state_size_felts
padding = (inp + keccak_func(inp)) * _block_size
segments.write_arg(ids.keccak_ptr_end, padding)`

const CAIRO_KECCAK_FINALIZE_V2 = `# Add dummy pairs of input and output.
_keccak_state_size_felts = int(ids.KECCAK_STATE_SIZE_FELTS)
_block_size = int(ids.BLOCK_SIZE)
assert 0 <= _keccak_state_size_felts < 100
assert 0 <= _block_size < 1000
inp = [0] * _keccak_state_size_felts
padding = (inp + keccak_func(inp)) * _block_size
segments.write_arg(ids.keccak_ptr_end, padding)`
