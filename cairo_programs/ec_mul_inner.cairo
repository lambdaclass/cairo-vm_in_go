%builtins range_check

from starkware.cairo.common.cairo_secp.ec import (
    EcPoint,
    ec_mul_inner,
)
from starkware.cairo.common.cairo_secp.bigint import BigInt3

func main{range_check_ptr: felt}() {
    // ec_mul_inner
    let (pow2, res) = ec_mul_inner(
        EcPoint(
            BigInt3(65162296, 359657, 04862662171381), BigInt3(-5166641367474701, -63029418, 793)
        ),
        123,
        298,
    );
    assert pow2 = EcPoint(
        BigInt3(30016796425722798916160189, 75045389156830800234717485, 13862403786096360935413684),
        BigInt3(43820690643633544357415586, 29808113745001228006676979, 15112469502208690731782390),
    );
    return ();
}
