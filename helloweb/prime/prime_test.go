package prime

import (
    "fmt"
    "testing"
    "time"
)

// maxEvalDuration is the acceptance-criteria budget (issue #53) for a
// single IsPrime evaluation near the 64-bit boundary.
const maxEvalDuration = 250 * time.Millisecond

func TestIsPrimeSmallPrimes(t *testing.T) {
    for _, n := range []int64{2, 3, 5, 7, 11, 23, 29, 37, 97, 107} {
        if !IsPrime(n) {
            t.Errorf("IsPrime(%d) = false, want true", n)
        }
    }
}

func TestIsPrimeSmallComposites(t *testing.T) {
    for _, n := range []int64{1, 4, 6, 8, 9, 10, 24, 25, 38, 99, 115} {
        if IsPrime(n) {
            t.Errorf("IsPrime(%d) = true, want false", n)
        }
    }
}

func TestIsPrimeLargerPrimes(t *testing.T) {
    for _, n := range []int64{7919, 104729, 15485863, 999999937} {
        if !IsPrime(n) {
            t.Errorf("IsPrime(%d) = false, want true", n)
        }
    }
}

// TestIsPrimeMaxInt64 covers 2^63-1, the maximum signed 64-bit integer,
// which is not prime.
func TestIsPrimeMaxInt64(t *testing.T) {
    n := int64(9223372036854775807)
    if IsPrime(n) {
        t.Errorf("IsPrime(%d) = true, want false", n)
    }
}

// TestIsPrimeLargestInt64PrimePerf covers 2^63-25, the largest prime
// that fits in a signed 64-bit integer, and the issue #53 acceptance
// criterion that evaluating it takes 250ms or less.
func TestIsPrimeLargestInt64PrimePerf(t *testing.T) {
    n := int64(9223372036854775783)

    start := time.Now()
    got := IsPrime(n)
    elapsed := time.Since(start)

    if !got {
        t.Errorf("IsPrime(%d) = false, want true", n)
    }
    if elapsed > maxEvalDuration {
        t.Errorf("IsPrime(%d) took %v, want <= %v", n, elapsed, maxEvalDuration)
    }
}

// TestIsPrimeUpperBoundaryPerf covers the issue #53 acceptance criterion
// that every value from 2^63-24 up to 2^63-1 evaluates in 250ms or less.
// All 24 are composite (2^63-25 is the largest prime below 2^63).
func TestIsPrimeUpperBoundaryPerf(t *testing.T) {
    const maxInt64 = int64(9223372036854775807)
    // Computed as an offset from maxInt64 rather than a directly
    // incrementing loop variable: incrementing n up to and past maxInt64
    // would overflow back to math.MinInt64, which is still <= maxInt64,
    // turning this into a practically infinite loop.
    for offset := int64(23); offset >= 0; offset-- {
        n := maxInt64 - offset
        start := time.Now()
        got := IsPrime(n)
        elapsed := time.Since(start)

        if got {
            t.Errorf("IsPrime(%d) = true, want false", n)
        }
        if elapsed > maxEvalDuration {
            t.Errorf("IsPrime(%d) took %v, want <= %v", n, elapsed, maxEvalDuration)
        }
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
