%builtins range_check

from starkware.cairo.common.math import assert_nn

func main{range_check_ptr: felt}() {
    let x = 64;
    tempvar y = 64 * 64;
    assert_nn(1);
    assert_nn(x);
    assert_nn(y);
    return ();
}
