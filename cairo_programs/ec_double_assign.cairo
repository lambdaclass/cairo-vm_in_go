%builtins range_check
from starkware.cairo.common.cairo_secp.bigint import BigInt3
from starkware.cairo.common.cairo_secp.ec import ec_double, EcPoint

func main{range_check_ptr}() {
    let p = EcPoint(BigInt3(1,2,3), BigInt3(4,5,6));

    let (r) = ec_double(p);

    assert r.x.d0 = 15463639180909693576579425;
    assert r.x.d1 = 18412232947780787290771221;
    assert r.x.d2 = 2302636566907525872042731;

    assert r.y.d0 = 62720835442754730087165024;
    assert r.y.d1 = 51587896485732116326812460;
    assert r.y.d2 = 1463255073263285938516131;

    return ();
}
