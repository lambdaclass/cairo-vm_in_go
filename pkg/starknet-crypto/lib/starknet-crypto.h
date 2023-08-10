#include <stdint.h>

typedef uint64_t limb_t;

/* A 256 bit prime field element (felt), represented as four limbs (integers).
 */
typedef limb_t felt_t[4];

typedef felt_t poseidon_state_t[4];

void poseidon_permute(poseidon_state_t);
