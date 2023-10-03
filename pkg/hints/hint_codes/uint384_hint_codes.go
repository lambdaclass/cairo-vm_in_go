package hint_codes

const UINT384_SPLIT_128 = `ids.low = ids.a & ((1<<128) - 1)
ids.high = ids.a >> 128`
