%builtins range_check

from starkware.cairo.common.math import abs_value

func main{range_check_ptr: felt}() {
    assert abs_value(-1) = 1
    assert abs_value(17) = 17
    assert abs_value(-21938134) = 21938134
    return ();
}
