%builtins output

from starkware.cairo.common.math import assert_not_equal

func main{output_ptr: felt*}() {
    assert_not_equal(17, 7);
    assert_not_equal(cast(output_ptr, felt), cast(output_ptr + 1, felt));
    assert_not_equal(-1, 1);
    return ();
}