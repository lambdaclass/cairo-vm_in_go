package hints

const IS_NN = "memory[ap] = 0 if 0 <= (ids.a % PRIME) < range_check_builtin.bound else 1"

const IS_NN_OUT_OF_RANGE = "memory[ap] = 0 if 0 <= ((-ids.a - 1) % PRIME) < range_check_builtin.bound else 1"

const IS_LE_FELT = "memory[ap] = 0 if (ids.a % PRIME) <= (ids.b % PRIME) else 1"
