%builtins output pedersen range_check

from starkware.cairo.common.cairo_builtins import HashBuiltin
from starkware.cairo.common.hash import hash2

func main{output_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}() {    
    let (seed_2) = hash2{hash_ptr=pedersen_ptr}(234, 123897213);
    assert seed_2 = 2528904803005991377642213282618516374663011807602690623037041511517940816555;
    return ();
}
