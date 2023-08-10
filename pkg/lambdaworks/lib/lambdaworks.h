#include <stdint.h>

typedef uint64_t limb_t;

/* A 256 bit prime field element (felt), represented as four limbs (integers).
 */
typedef limb_t felt_t[4];

/* Gets a felt_t representing the "value" number, in montgomery format. */
void from(felt_t result, uint64_t value);

/*Gets a felt_t representing the "value" hexadecimal string, in montgomery
 * format. */
void from_hex(felt_t result, char *value);

/*Gets a felt_t representing the "value" decimal string, in montgomery format.
 */
void from_dec_str(felt_t result, char *value);

/* Converts a felt_t to bytes in little-endian representation. */
void to_le_bytes(uint8_t result[32], felt_t value);

/* Gets a felt_t representing 0 */
void zero(felt_t result);

/* Gets a felt_t representing 1 */
void one(felt_t result);

/* Writes the result variable with the sum of a and b felts. */
void add(felt_t a, felt_t b, felt_t result);

/* Writes the result variable with a - b. */
void sub(felt_t a, felt_t b, felt_t result);

/* Writes the result variable with a * b. */
void mul(felt_t a, felt_t b, felt_t result);

/* Writes the result variable with a / b. */
void lw_div(felt_t a, felt_t b, felt_t result);

/* Returns the minimum number of bits needed to represent the felt */
limb_t bits(felt_t a);

/* writes the result variable with a & b */
void felt_and(felt_t a, felt_t b, felt_t result);

/* writes the result variable with a | b */
void felt_or(felt_t a, felt_t b, felt_t result);

void felt_pow(felt_t a, felt_t b, felt_t result);