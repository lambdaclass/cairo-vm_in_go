%builtins range_check bitwise
from starkware.cairo.common.uint256 import (
    Uint256,
    uint256_sqrt,
)
from starkware.cairo.common.cairo_builtins import BitwiseBuiltin

func main{range_check_ptr: felt, bitwise_ptr: BitwiseBuiltin*}() {
    let n = Uint256(0, 157560248172239344387757911110183813120);
    let res = uint256_sqrt(n);
    return ();
}
