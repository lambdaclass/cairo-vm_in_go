package hint_codes

const ADD_SEGMENT = "memory[ap] = segments.add()"
const VM_EXIT_SCOPE = "vm_exit_scope()"
const VM_ENTER_SCOPE = "vm_enter_scope()"
const MEMCPY_ENTER_SCOPE = "vm_enter_scope({'n': ids.len})"
const MEMCPY_CONTINUE_COPYING = "n -= 1\nids.continue_copying = 1 if n > 0 else 0"
