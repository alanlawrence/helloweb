package series

import (
    "fmt"
    "strconv"
    "net/url"
)

// Globals
var DEBUG = false


// Inputs: a and d are floats. n is integer >= 1 (represented as float)
// Outputs the sum and error string (which is "" if all is good).
func ARSum(a, d, n float64) (arSum float64, err string) {

    // Check that the coeefs have valid values. Return an error string if not.
    err = ""
    if n < 1 {
        err += fmt.Sprintf("Expected n >= 1 but got %v", n)
    }

    arSum = 0 // Default return value
    // If coeffs are good, then compute the sum.
    if err == "" {
        // Improve precision by increasing the integer part. But only
        // if we need to. We could improve this crude approach by making
        // sure we always have at least 6 sf in the integer part for the
        // computation.
        shift := 1.0
        // If there are any decimals in a or d, we'll shift up first.
        if (a != float64(int64(a)) && d != float64(int64(d))) {
            shift = 1e6
        }
        a = a * shift
        d = d * shift
        arSum = (n / 2.0) * (2.0 * a + (n - 1.0) * d) / shift
        if DEBUG { fmt.Printf("arSum=%v\n", arSum) }
    }

    return arSum, err
}

func GenerateHtml(coeff_a, coeff_d, coeff_n float64, arSum float64,
                  err string) (htmlStr string) {

    htmlStr += fmt.Sprintf("For a=%v, d=%v, n=%v, the sum is %v",
                                coeff_a, coeff_d, coeff_n, arSum)
    if err != "" {
        htmlStr += fmt.Sprintf("<br>Oops! %v", err)
    }

    return htmlStr
}

// Computes the final html str result from the query values.
func ARCalc(urlStr string, v url.Values) (htmlStr string) {
    // Extract args one by one and compute the binary signature.
    // Return the binary signature as the sum for development purposes.
    // Search for each parameter and compute the binary signature.
    binSig   := 0.0
    coeffA   := 0.0
    coeffD   := 0.0
    coeffN   := 0.0
    coeffs, ok := v["a"]
    if ok {
        binSig += 1
        coeffAStr := coeffs[0]
        coeffA, _ = strconv.ParseFloat(coeffAStr, 64)
    }
    coeffs, ok = v["d"]
    if ok {
        binSig += 2
        coeffDStr := coeffs[0]
        coeffD, _ = strconv.ParseFloat(coeffDStr, 64)
    }
    coeffs, ok = v["n"]
    if ok {
        binSig += 4
        coeffNStr := coeffs[0]
        coeffN, _ = strconv.ParseFloat(coeffNStr, 64)
    }

    if binSig == 7 {
        arSum, err := ARSum(coeffA, coeffD, coeffN)
        htmlStr = GenerateHtml(coeffA, coeffD, coeffN, arSum, err)
    } else {
        htmlStr = fmt.Sprintf("Expected 3 params /a, d, n/ but got %v", urlStr)
    }
    return htmlStr
}
