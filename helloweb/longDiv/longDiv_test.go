package longDiv

import "testing"

func TestLongDiv(t *testing.T) {

    if !SelfTest() {
        t.Errorf("LongDiv self test failed")
    }

    // Test PrintWorking. Note the stdout is not printed
    // unless there is a test failure, i.e. t.Errorf(...)
    PrintWorking(LongDiv(12, 143))
}

// TestLongDivNumSmallerThanDenom exercises the case where the numerator
// loop never executes because num < denom, leaving a zero quotient and
// the numerator as the remainder.
func TestLongDivNumSmallerThanDenom(t *testing.T) {
    _, quotient, interim, rem := LongDiv(5, 3)

    if quotient.Sprint() != "0" {
        t.Errorf("Expected quotient 0 but got %v", quotient.Sprint())
    }
    if rem != 3 {
        t.Errorf("Expected remainder 3 but got %v", rem)
    }
    if len(interim) != 3 {
        t.Errorf("Expected len(interim)=3 but got %v", len(interim))
    }
}

// TestLongDivMultiDigitQuotient exercises a multi-digit quotient with an
// exact division and no remainder (900 / 4 = 225).
func TestLongDivMultiDigitQuotient(t *testing.T) {
    _, quotient, _, rem := LongDiv(4, 900)

    if quotient.Sprint() != "225" {
        t.Errorf("Expected quotient 225 but got %v", quotient.Sprint())
    }
    if rem != 0 {
        t.Errorf("Expected remainder 0 but got %v", rem)
    }
}

// TestLongDivIncreasePrecisionBranch exercises the branch where the
// denominator does not divide the leading digit(s) of the numerator, so
// the interim numerator has to be widened before division succeeds
// (100 / 23 = 4 remainder 8).
func TestLongDivIncreasePrecisionBranch(t *testing.T) {
    _, quotient, _, rem := LongDiv(23, 100)

    if quotient.Sprint() != "  4" {
        t.Errorf("Expected quotient \"  4\" but got %q", quotient.Sprint())
    }
    if rem != 8 {
        t.Errorf("Expected remainder 8 but got %v", rem)
    }
}

// TestLongDivInvalidDenom exercises the guard for a denominator less
// than 1, which returns early without panicking.
func TestLongDivInvalidDenom(t *testing.T) {
    digitsD, quotient, _, rem := LongDiv(0, 5)

    if digitsD.Sprint() != "0" {
        t.Errorf("Expected digitsD \"0\" but got %q", digitsD.Sprint())
    }
    if quotient.Sprint() != "0" {
        t.Errorf("Expected quotient \"0\" but got %q", quotient.Sprint())
    }
    if rem != 5 {
        t.Errorf("Expected remainder 5 but got %v", rem)
    }
}

func TestGenerateHtml(t *testing.T) {

    cases := []struct {
        name     string
        denom    int
        num      int
        wantHtml string
    }{
        {
            "single digit division into double digit numerator with remainder",
            7, 15,
            "   2 remainder=1<br> .--<br>7|15<br>  14<br>  --<br>   1<br>",
        },
        {
            "denom exactly divides",
            6, 6,
            "  1 remainder=0<br> .-<br>6|6<br>  6<br>  -<br>  0<br>",
        },
        {
            "numerator smaller than denominator",
            5, 3,
            "  0 remainder=3<br> .-<br>5|3<br>",
        },
        {
            "multi-digit quotient, exact division",
            4, 900,
            "  225 remainder=0<br> .---<br>4|900<br>  800<br>  ---<br>" +
                "  100<br>   80<br>  ---<br>   20<br>   20<br>  ---<br>    0<br>",
        },
        {
            "double digit division into triple digit numerator with remainder",
            23, 100,
            "     4 remainder=8<br>  .---<br>23|100<br>    92<br>   ---<br>     8<br>",
        },
    }

    for _, c := range cases {
        // Also exercise PrintWorking's stdout path for each case;
        // stdout is only printed on failure of a later assertion.
        PrintWorking(LongDiv(c.denom, c.num))

        htmlStr := GenerateHtml(LongDiv(c.denom, c.num))
        if htmlStr != c.wantHtml {
            t.Errorf("%v: expected html=\n%q\nbut got\n%q",
                                c.name,          c.wantHtml,  htmlStr)
        }
    }
}

func TestCalculateHtml(t *testing.T) {
    validCases := []struct {
        name     string
        denom    float64
        num      float64
        wantHtml string
    }{
        {
            "positive integers",
            7, 15,
            "   2 remainder=1<br> .--<br>7|15<br>  14<br>  --<br>   1<br>",
        },
    }
    for _, c := range validCases {
        got := CalculateHtml(c.denom, c.num)
        if got != c.wantHtml {
            t.Errorf("%v: expected html=\n%q\nbut got\n%q",
                                c.name,          c.wantHtml,  got)
        }
    }

    invalidCases := []struct {
        name  string
        denom float64
        num   float64
    }{
        {"decimal denominator", 3.5, 4},
        {"decimal numerator", 3, 4.5},
        {"negative denominator", -3, 4},
        {"negative numerator", 3, -4},
        {"zero denominator", 0, 4},
        {"zero numerator", 3, 0},
    }
    for _, c := range invalidCases {
        got := CalculateHtml(c.denom, c.num)
        if got != invalidInputMsg {
            t.Errorf("%v: expected error message %q but got %q",
                                c.name,          invalidInputMsg, got)
        }
    }
}
