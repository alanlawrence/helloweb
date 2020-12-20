package digits

import (
    "fmt"
    "math"
)

// package globals
var DEBUG = false

// Digits class.
type Digits struct {
    digits []float64
}

// Initialise to the length of num and populate the digits
// from num.
func (d *Digits) Init(num int) () {
    lenN := 1
    if num > 0 {
        lenN = int(math.Trunc(math.Log10(float64(num)) + 1.0))
    }
    d.digits = make([]float64,lenN)
    base := 10.0
    for di := 0; di < lenN; di++ {
        d.digits[di] = math.Mod(
            math.Trunc(float64(num)/math.Pow10(lenN-1-di)), base)
    }
}

// Initialise to the fixed length given and populate the digits
// with the supplied number.
// If the number has less digits than length, then left pad with zeros
// If the number has more digits than length, then truncate on the left
func (d *Digits) InitFixed(lenN int, num int) () {
    d.digits = make([]float64,lenN)
    base := 10.0
    for di := 0; di < lenN; di++ {
        d.digits[di] = math.Mod(
            math.Trunc(float64(num)/math.Pow10(lenN-1-di)), base)
    }
}

func (d *Digits) SetDigit(di int, digit float64) () {
    if (di >= 0 && di < len(d.digits)) {
        d.digits[di] = digit
    }
    // Else no op
}

func (d Digits) Digit(di int) (digit float64) {
    if (di >= 0 && di < len(d.digits)) {
        digit = d.digits[di]
    }
    // Else no op
    return digit
}

func (d *Digits) Print(name string) () {
    fmt.Printf("Digits of %v: ", name)
    for i := range d.digits {
        fmt.Printf("%v", d.digits[i])
    }
    fmt.Printf("\n")
}

func (d *Digits) Sprint() (digits string) {

    leading := true
    for i := range d.digits {
        // Don't print leading zeros unless zero is the only digit.
        if leading && len(d.digits)>1 && d.digits[i] == 0 {
            digits += " "
        } else {
            digits += fmt.Sprintf("%v", d.digits[i])
            leading = false
        }
    }
    return digits
}

func (d *Digits) SprintDigit(di int) (digit string) {

    if (di >= 0 && di < len(d.digits)) {
        digit = fmt.Sprintf("%v", d.digits[di])
    }
    return digit
}

func (d *Digits) PartN(j int) (partN float64) {
    partN = 0
    for k := 0; k <= j; k++ {
        // Compute the part of numerator length as powers of 10 of
        // the digits in the appropriate interim results row.
        partN += d.digits[k] * math.Pow10(j-k)
        if DEBUG {fmt.Printf("partN=%v, j=%v, k=%v\n",
                              partN,    j,    k)}
    }
    return partN
}

func (d *Digits) Len() (lend int) {
    lend = len(d.digits)
    return lend
}

func (d *Digits) Compare(d2 Digits) (equal bool) {
    // First compare lengths
    equal = (len(d.digits) == len(d2.digits))
    if !equal {
        fmt.Printf("Expected d=%v but got d2=%v\n", d.Sprint(), d2.Sprint())
        fmt.Printf("Expected len(d)=%v but got len(d2)=%v\n",
                             len(d.digits),    len(d2.digits))
        return equal
    }

    // Then compare digit by digit
    for i := 0; i <= len(d.digits) - 1; i++ {
        equal = d.digits[i] == d2.digits[i]
        if !equal {
            fmt.Printf("Expected d=%v but got d2=%v\n", d.Sprint(), d2.Sprint())
            fmt.Printf("Expected d[%v]=%v but got d2[%v]=%v\n",
                        i, d.digits[i], i, d2.digits[i])
            break;
        }
    }

    return equal
}
