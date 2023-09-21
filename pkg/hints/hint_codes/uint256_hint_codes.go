package hint_codes

const UINT256_ADD = "sum_low = ids.a.low + ids.b.low \n ids.carry_low = 1 if sum_low >= ids.SHIFT else 0 \n sum_high = ids.a.high + ids.b.high + ids.carry_low \n ids.carry_high = 1 if sum_high >= ids.SHIFT else 0"
