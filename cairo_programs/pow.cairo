%builtins range_check

from starkware.cairo.common.pow import pow

func main{range_check_ptr: felt}() {
    let (x) = pow(2, 3);
    assert x = 8;
    let (y) = pow(10, 6);
    assert y = 1000000;
    let (z) = pow(152, 25);
    assert z = 3516330588649452857943715400722794159857838650852114432;
    let (u) = pow(-2, 3);
    assert (u) = -8;
    let (v) = pow(-25, 31);
    assert (v) = -21684043449710088680149056017398834228515625;

    return ();
}
