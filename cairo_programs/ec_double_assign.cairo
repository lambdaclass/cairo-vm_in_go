%builtins range_check
from starkware.cairo.common.cairo_secp.bigint import BigInt3, nondet_bigint3
struct EcPoint {
    x: BigInt3,
    y: BigInt3,
}

func ec_double{range_check_ptr}(point: EcPoint, slope: BigInt3) -> (res: BigInt3) {
    %{
        from starkware.cairo.common.cairo_secp.secp_utils import pack
        SECP_P = 2**255-19

        slope = pack(ids.slope, PRIME)
        x = pack(ids.point.x, PRIME)
        y = pack(ids.point.y, PRIME)

        value = new_x = (pow(slope, 2, SECP_P) - 2 * x) % SECP_P
    %}

    //let (new_x: BigInt3) = nondet_bigint3();
    return (res=slope);
}

func main{range_check_ptr}() {
    let p = EcPoint(BigInt3(0,0,2), BigInt3(4,5,6));
    let s = BigInt3(0,0,3);
    let (res) = ec_double(p, s);
    assert res.d0 = 5;
    assert res.d1 = 0;
    assert res.d2 = 0;
    return ();
}
