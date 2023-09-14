#include <stdint.h>
#include <stddef.h> 
#include <stdlib.h>

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

/* Converts a felt_t to bytes in big-endian representation. */
void to_be_bytes(uint8_t result[32], felt_t value);

/* Converts a felt_t to a String representation. */
void to_hex_string(char *string, felt_t value);

/* Converts an array of bytes in little-endian representation to a felt_t. */
void from_le_bytes(felt_t result, uint8_t bytes[32]);

/* Converts an array of bytes in big-endian representation to a felt_t. */
void from_be_bytes(felt_t result, uint8_t bytes[32]);

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

/* writes the result variable with a ^ b */
void felt_xor(felt_t a, felt_t b, felt_t result);

/* writes the result variable with a << num */
void felt_shl(felt_t a, uint64_t num, felt_t result);

/* writes the result variable with a.pow(num) */
void felt_pow_uint(felt_t a, uint32_t num, felt_t result);

/* returns the representation of a felt to string */
char* to_signed_felt(felt_t value);

/* frees a pointer to a string */
void free_string(char* ptr);

/* writes the result variable with a >> num */
void felt_shr(felt_t a, size_t b, felt_t result);

/*
Compares x and y and returns:

	-1 if a <  b
	 0 if a == b
	+1 if a >  b
*/
int cmp(felt_t a, felt_t b);
