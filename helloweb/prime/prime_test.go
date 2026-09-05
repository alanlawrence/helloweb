package prime

import (
    "fmt"
    "testing"
)

func TestIsPrimeSmallPrimes(t *testing.T) {
    for _, n := range []int64{1, 2, 3, 5, 7, 11, 23, 29, 37, 97, 107} {
        if !IsPrime(n) {
            t.Errorf("IsPrime(%d) = false, want true", n)
        }
    }
}

func TestIsPrimeSmallComposites(t *testing.T) {
    for _, n := range []int64{4, 6, 8, 9, 10, 24, 25, 38, 99, 115} {
        if IsPrime(n) {
            t.Errorf("IsPrime(%d) = true, want false", n)
        }
    }
}

// TestIsPrimeLargerPrimes exercises the 6k+/-1 loop with primes large
// enough to require it, but small enough to stay fast: this algorithm is
// O(sqrt(n)) trial division, which is impractical near the 64-bit
// boundary (n ~ 2^63 means sqrt(n) ~ 3 billion loop iterations). Proving
// primality of numbers at that scale is deferred to issue #53, which
// will replace this with a deterministic Miller-Rabin test.
func TestIsPrimeLargerPrimes(t *testing.T) {
    for _, n := range []int64{7919, 104729, 15485863, 999999937} {
        if !IsPrime(n) {
            t.Errorf("IsPrime(%d) = false, want true", n)
        }
    }
}

// TestIsPrimeMaxInt64 covers 2^63-1, the maximum signed 64-bit integer.
// It is not prime and, usefully, has a small factor, so trial division
// resolves it quickly despite its size.
func TestIsPrimeMaxInt64(t *testing.T) {
    n := int64(9223372036854775807)
    if IsPrime(n) {
        t.Errorf("IsPrime(%d) = true, want false", n)
    }
}

func TestCalculateHtmlValidNumbers(t *testing.T) {
    tests := []struct {
        numberStr string
        want      string
    }{
        {"2", "2 is prime!"},
        {"4", "4 is not prime :-("},
        {"13", "13 is prime!"},
        {"15485863", "15485863 is prime!"},
        {"999999937", "999999937 is prime!"},
    }
    for _, tc := range tests {
        if got := CalculateHtml(tc.numberStr); got != tc.want {
            t.Errorf("CalculateHtml(%q) = %q, want %q", tc.numberStr, got, tc.want)
        }
    }
}

// TestCalculateHtmlOverflowClampsToMaxInt64 covers the acceptance
// criteria requiring the 12th Mersenne prime (2^127-1), 2^63, and
// 2^63+1 to all be treated as 2^63-1, which is not prime. These stay
// fast because the clamped value (2^63-1) has a small factor, unlike an
// arbitrary number of that magnitude.
func TestCalculateHtmlOverflowClampsToMaxInt64(t *testing.T) {
    want := fmt.Sprintf("%v %v", int64(9223372036854775807), "is not prime :-(")
    tests := []string{
        "170141183460469231731687303715884105727", // 12th Mersenne prime, 2^127-1
        "9223372036854775808",                     // 2^63
        "9223372036854775809",                     // 2^63+1
    }
    for _, numberStr := range tests {
        if got := CalculateHtml(numberStr); got != want {
            t.Errorf("CalculateHtml(%q) = %q, want %q", numberStr, got, want)
        }
    }
}

func TestCalculateHtmlInvalidInput(t *testing.T) {
    for _, numberStr := range []string{"-5", "0", "-9223372036854775809", "abc", "3.14", ""} {
        if got := CalculateHtml(numberStr); got != invalidInputMsg {
            t.Errorf("CalculateHtml(%q) = %q, want %q", numberStr, got, invalidInputMsg)
        }
    }
}
