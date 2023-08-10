#include <stdint.h>

typedef uint8_t byte_t;

/* A 256 bit prime field element (felt), represented as four limbs (integers).
 */
typedef byte_t felt_t[32];


void poseidon_permute(felt_t, felt_t, felt_t);
