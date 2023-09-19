package hints

const MEMSET_ENTER_SCOPE = "vm_enter_scope({'n': ids.n})"
const MEMSET_CONTINUE_LOOP = "n -= 1\nids.continue_loop = 1 if n > 0 else 0"
