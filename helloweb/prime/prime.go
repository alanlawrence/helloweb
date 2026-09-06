// Package prime provides primality testing and the HTML fragment
// generation used by the /prime webserver endpoint.
package prime

import (
    "fmt"
    "math/bits"
    "strconv"
)

const invalidInputMsg = "Only positive integers are valid inputs"

// Bases forming a deterministic Miller-Rabin witness set for every
// n < 2^64. See https://miller-rabin.appspot.com/ for the derivation.
var millerRabinBases = []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}

// modMul calculates (a * b) % mod without 64-bit overflow, using the
// full 128-bit product.
func modMul(a, b, mod uint64) uint64 {
    hi, lo := bits.Mul64(a, b)
    _, rem := bits.Div64(hi, lo, mod)
    return rem
}

// modPow calculates (base^exp) % mod.
func modPow(base, exp, mod uint64) uint64 {
    result := uint64(1)
    base = base % mod
    for exp > 0 {
        if exp&1 == 1 {
            result = modMul(result, base, mod)
        }
        base = modMul(base, base, mod)
        exp >>= 1
    }
    return result
}

// IsPrime reports whether n is a prime number, using the deterministic
// Miller-Rabin test with a fixed witness set valid for all n < 2^64.
func IsPrime(n int64) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    un := uint64(n)
    if un%2 == 0 || un%3 == 0 {
        return false
    }

    // Factor n - 1 into 2^s * d.
    d := un - 1
    s := 0
    for d%2 == 0 {
        d /= 2
        s++
    }

    for _, a := range millerRabinBases {
        if a >= un {
            a %= un
            if a == 0 {
                continue
            }
        }

        x := modPow(a, d, un)
        if x == 1 || x == un-1 {
            continue
        }

        composite := true
        for r := 0; r < s-1; r++ {
            x = modMul(x, x, un)
            if x == un-1 {
                composite = false
                break
            }
        }
        if composite {
            return false
        }
    }

    return true
}

// GenerateHtml formats the primality result for number as the HTML
// fragment the /prime endpoint responds with.
func GenerateHtml(number int64, isPrime bool) string {
    result := "is not prime :-("
    if isPrime {
        result = "is prime!"
    }
    return fmt.Sprintf("%v %v", number, result)
}

// CalculateHtml parses numberStr as a positive integer and returns the
// primality result as an HTML fragment. Integers outside the signed
// 64-bit range are clamped to the maximum magnitude of that range (e.g.
// 2^63 and above are treated as 2^63-1) rather than rejected. Anything
// that isn't a positive integer (malformed input, negative, or zero)
// returns an HTML fragment containing an error message instead.
func CalculateHtml(numberStr string) string {
    const base int = 10
    const bitSize int = 64
    number, err := strconv.ParseInt(numberStr, base, bitSize)
    if err != nil {
        numErr, ok := err.(*strconv.NumError)
        if !ok || numErr.Err != strconv.ErrRange {
            return invalidInputMsg
        }
        // Out of range: ParseInt already returned the maximum magnitude
        // int64 value of the appropriate sign, which we use as-is.
    }
    if number <= 0 {
        return invalidInputMsg
    }
    return GenerateHtml(number, IsPrime(number))
}
