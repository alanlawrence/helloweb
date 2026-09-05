// Package prime provides primality testing and the HTML fragment
// generation used by the /prime webserver endpoint.
package prime

import (
    "fmt"
    "strconv"
)

const invalidInputMsg = "Only positive integers are valid inputs"

// IsPrime reports whether n is a prime number.
func IsPrime(n int64) bool {

    isPrime := true
    done := false
    if n <= 3 {
        isPrime = n >= 1
        done = true
    } else if (n % 2 == 0 || n % 3 == 0) {
        isPrime = false
        done = true
    }

    // All the numbers below 25 are either divisible by 2 or 3 (tested above)
    // or are prime.
    if (!done && n < 25) {
        isPrime = true
        done = true
    }

    // Now exploit the property that all primes >=6 are of the form 6k+1 or 6k-1
    // since 2 divides 6k, 6k+2 and 6k+4, and 3 divides 6k+3
    // which leaves 6k+1 and 6k+5 (== 6k'-1, where k'=k+1).

    // So we test all numbers of the form 6k+/-1 such that
    //     6k+/-1     <= sqrt(n)
    // ==> (6k+/-1)^2 <= n

    i := int64(5)
    // This generates the pair 5 and 7 for the first iteration, k = 1
    for (!done && i*i <= n) {

        if n % i == 0 {
            isPrime = false
            done = true
        } else if n % (i+2) == 0 {
            isPrime = false
            done = true
        }
        // Advance to next iteration. Imagine k += 1
        i += 6
    }
    // else must be prime hence return isPrime default of true.

    return isPrime
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
    number, err := strconv.ParseInt(numberStr, 10, 64)
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
