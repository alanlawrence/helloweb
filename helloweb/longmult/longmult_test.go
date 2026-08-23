package longmult

import (
    "reflect"
    "testing"
)

func TestLongMultSingleDigit(t *testing.T) {
    digitsA, digitsB, product, carrys := LongMult(9, 9)

    wantA := []float64{9}
    wantB := []float64{9}
    wantProduct := [][]float64{{8, 1}, {8, 1}}
    wantCarrys := [][]float64{{0, 8}}

    if !reflect.DeepEqual(digitsA, wantA) {
        t.Errorf("digitsA: expected %v but got %v", wantA, digitsA)
    }
    if !reflect.DeepEqual(digitsB, wantB) {
        t.Errorf("digitsB: expected %v but got %v", wantB, digitsB)
    }
    if !reflect.DeepEqual(product, wantProduct) {
        t.Errorf("product: expected %v but got %v", wantProduct, product)
    }
    if !reflect.DeepEqual(carrys, wantCarrys) {
        t.Errorf("carrys: expected %v but got %v", wantCarrys, carrys)
    }
}

// TestLongMultMultiDigit exercises the carry logic across multiple digits
// and multiple rows of the multiplier (24 x 418, i.e. bnum=24, anum=418).
func TestLongMultMultiDigit(t *testing.T) {
    digitsA, digitsB, product, carrys := LongMult(24, 418)

    wantA := []float64{4, 1, 8}
    wantB := []float64{2, 4}
    wantProduct := [][]float64{
        {0, 8, 3, 6, 0},
        {0, 1, 6, 7, 2},
        {1, 0, 0, 3, 2},
    }
    wantCarrys := [][]float64{
        {0, 0, 0, 1, 0},
        {0, 0, 1, 0, 3},
    }

    if !reflect.DeepEqual(digitsA, wantA) {
        t.Errorf("digitsA: expected %v but got %v", wantA, digitsA)
    }
    if !reflect.DeepEqual(digitsB, wantB) {
        t.Errorf("digitsB: expected %v but got %v", wantB, digitsB)
    }
    if !reflect.DeepEqual(product, wantProduct) {
        t.Errorf("product: expected %v but got %v", wantProduct, product)
    }
    if !reflect.DeepEqual(carrys, wantCarrys) {
        t.Errorf("carrys: expected %v but got %v", wantCarrys, carrys)
    }

    // The last row of product is the total, read left to right.
    wantTotal := []float64{1, 0, 0, 3, 2} // 10032 == 418 * 24
    gotTotal := product[len(product)-1]
    if !reflect.DeepEqual(gotTotal, wantTotal) {
        t.Errorf("total: expected %v but got %v", wantTotal, gotTotal)
    }
}

// TestLongMultRepeatedCarryDigit exercises a case (99 x 99) where the
// carry from one digit multiplication feeds into the next, and the carry
// total row itself also carries.
func TestLongMultRepeatedCarryDigit(t *testing.T) {
    _, _, product, _ := LongMult(99, 99)

    wantTotal := []float64{9, 8, 0, 1} // 9801 == 99 * 99
    gotTotal := product[len(product)-1]
    if !reflect.DeepEqual(gotTotal, wantTotal) {
        t.Errorf("total: expected %v but got %v", wantTotal, gotTotal)
    }
}

func TestPrintWorking(t *testing.T) {
    // Not asserted on: stdout is only useful for eyeballing the working.
    // This exercises PrintWorking so it is covered.
    PrintWorking(LongMult(24, 418))
}

func TestGenerateHtml(t *testing.T) {
    cases := []struct {
        name    string
        bnum    int
        anum    int
        wantHtml string
    }{
        {
            "single digit",
            7, 8,
            "            [8]<br>          x [7]<br><br>" +
                "          5 6 <br>carry->   5 <br>Total     5 6 <br>",
        },
        {
            "multi digit with carrying",
            24, 418,
            "            [4 1 8]<br>            x [2 4]<br><br>" +
                "            1 6 7 2 <br>carry->     1   3 <br>" +
                "            8 3 6 0 <br>carry->       1   <br>" +
                "Total     1 0 0 3 2 <br>",
        },
    }

    for _, c := range cases {
        got := GenerateHtml(LongMult(c.bnum, c.anum))
        if got != c.wantHtml {
            t.Errorf("%v: expected html=\n%q\nbut got\n%q",
                                c.name,          c.wantHtml,  got)
        }
    }
}

func TestCalculateHtml(t *testing.T) {
    validCases := []struct {
        name     string
        a        float64
        b        float64
        wantHtml string
    }{
        {
            "positive integers",
            8, 7,
            "            [8]<br>          x [7]<br><br>" +
                "          5 6 <br>carry->   5 <br>Total     5 6 <br>",
        },
    }
    for _, c := range validCases {
        got := CalculateHtml(c.a, c.b)
        if got != c.wantHtml {
            t.Errorf("%v: expected html=\n%q\nbut got\n%q",
                                c.name,          c.wantHtml,  got)
        }
    }

    invalidCases := []struct {
        name string
        a    float64
        b    float64
    }{
        {"decimal first operand", 3.5, 4},
        {"decimal second operand", 3, 4.5},
        {"negative first operand", -3, 4},
        {"negative second operand", 3, -4},
        {"zero first operand", 0, 4},
        {"zero second operand", 3, 0},
    }
    for _, c := range invalidCases {
        got := CalculateHtml(c.a, c.b)
        if got != invalidInputMsg {
            t.Errorf("%v: expected error message %q but got %q",
                                c.name,          invalidInputMsg, got)
        }
    }
}
