package longDiv

import (
    "fmt"
    "math"
    "strings"
    d "helloweb/digits"
)

// Globals
var DEBUG = false


// Notes:
// Separates out interim quotient from quotient
// to deal with multi digit interim numerators
// e.g
//     011
//  12 143 interim quotient 010
//     120  
//      23 interim quotient  01
//      12 
//      11 
// final quotient becomes: addition of interim quotients
func LongDiv(denom int, num int) (digitsD d.Digits, quotient d.Digits,
                                         interim []d.Digits, numr int) {

    var digitsN d.Digits
    digitsN.Init(num)
    if DEBUG {digitsN.Print("num")}

    digitsD.Init(denom)
    if DEBUG {digitsD.Print("denom")}

    //quotient := array [lenN]
    //The quotient can never be longer than the numerator for integer division.
    lenN := digitsN.Len()
    quotient.InitFixed(lenN, 0)
    interimQ := make([]d.Digits, lenN)

    // interim := array [2 * lenN +1][lenN]
    //    Holds inital numerator then
    //    a pair of rows of interim results for each digit in the numerator.
    //        One for multiple of divisor, one for the interim remainder
    interim = make([]d.Digits, (2 * lenN) + 1)
    interim[0].Init(num)
    if DEBUG {interim[0].Print("interim[0]")}

    if denom < 1 {
        fmt.Printf("Error denom > 0 expected but got %v.\n", denom)
        return digitsD, quotient, interim, num
    }


    if DEBUG {
        fmt.Printf("%v", denom)
        fmt.Printf("\n")
        fmt.Printf("%v", num)
        fmt.Printf("\nQuotient: ")
    }

    i := 0 // Indexes interim quotient array horizontally.
    j := 0 // Indexes interim row horizontally which holds partN digits.
    partN := 0.0 // Part of numerator left after subtrating quotient * denom.
    interim0 := 0 // Holds interim denom * interim quotient calc row.
    irow := 0 // Indexes interim array vertically for each interim calculation.
    iq := 0 // Indexes quotient array vertically for each interim calculation.

    // Exit for loop when numerator is too small to be divided by denominator.
    for num >= denom {
        if DEBUG {fmt.Printf("------------partN loop %v---------------\n", i)}
        partN = interim[irow].PartN(j)
        // The interim quotient maps to the interim row containing 
        // the part of the numerator left. Hence half as many rows.
        iq = irow/2
        lenN = int(math.Trunc(math.Log10(float64(num)) + 1.0))
        // Make the interimQ row to the length required.
        interimQ[iq].InitFixed(lenN, 0)
        if DEBUG {fmt.Printf("interimQ[%v], len=%v\n", iq, interimQ[iq].Len())}
        // Assign the whole number of times denom divides the part of the
        // numerator left.
        interimQ[iq].SetDigit(i, math.Trunc(partN/float64(denom)))
        if DEBUG {
            fmt.Printf("interimQ[%v][%v]=%v\n",
                                  iq,  i,  interimQ[iq].SprintDigit(i))
        }
        if interimQ[iq].Digit(i) == 0 {
            // denom doesn't divide partN so increase length of partN
            i++
            j++
            if DEBUG {fmt.Printf("Increased j to %v\n", j)}
        } else {
            // denom does divide so compute new num as
            // num = num - quotient * denom
            interim0 = int(interimQ[iq].Digit(i) * float64(denom) *
                                             math.Pow10(lenN-i-1))
            if DEBUG {fmt.Printf("interim0=%v\n", interim0)}
            interim[irow+1].Init(interim0)
            if DEBUG {interim[irow+1].Print(
                            fmt.Sprintf("interim[%v]", irow+1))}
            num = num - interim0
            if DEBUG {fmt.Printf("num=%v\n", num)}
            // Now populate the interim row with the part of the numerator
            // left. This becomes partN as we work throught digits until
            // denom divides or num < denom.
            interim[irow+2].Init(num)
            if DEBUG {interim[irow+2].Print(
                                    fmt.Sprintf("interim[%v]", irow+2))}
            if DEBUG {interimQ[iq].Print("interimQ")}
            i = 0
            j = 0
            irow+=2
        }
    }
    if DEBUG {
        fmt.Printf("interim=%v\n", interim)
        fmt.Printf("interimQ=%v\n", interimQ)
    }

    // Overlay interim quotients onto quotient
    offset := 0
    for i := range interimQ {
        if DEBUG {fmt.Printf("i=%v\n", i)}
        for j:=interimQ[i].Len()-1; j >=0; j-- {
            if DEBUG {fmt.Printf("j=%v\n", j)}
            offset = quotient.Len() - interimQ[i].Len()
            if DEBUG {
                fmt.Printf("quotient[%v]=%v += interimQ[%v][%v]=%v\n",
                    j + offset, quotient.Digit(j), i, j, interimQ[i].Digit(j))
            }
            quotient.SetDigit(j + offset,
                            quotient.Digit(j + offset) + interimQ[i].Digit(j))
        }
    }
    if DEBUG {fmt.Printf("quotient=%v, remainder=%v\n", quotient, num)}

    return digitsD, quotient, interim, num
}

// TODO: guard pathological cases, e.g. interim is empty
func PrintWorking(digitsD d.Digits, quotient d.Digits,
                                     interim []d.Digits, rem int) () {
    // Print the quotient, a space and then remainder = the remainder
    // Print the divide bar
    // Print the denom and a vertical bar then the first row of interim
    // Print the interim rows
    // Get the padding right!

    fmt.Printf("==================================\n\n")
    lenD := digitsD.Len()
    fmtStr := fmt.Sprintf("%%%vv", lenD + 1 + quotient.Len())
    fmt.Printf(fmtStr + " remainder=%v\n", quotient.Sprint(), rem)
    lenI0 := interim[0].Len()
    fmt.Printf(fmtStr + "\n", strings.Repeat("-",lenI0 + 1))
    fmt.Printf("%v|%v\n", digitsD.Sprint(), interim[0].Sprint())
    endStop := lenD + 1 + lenI0
    for i := 1; i < len(interim); i++ {
        if interim[i].Len() < 1 {
            break;
        }
        fmtStr = fmt.Sprintf("%%%vv", endStop)
        fmt.Printf(fmtStr + "\n", interim[i].Sprint())
        if math.Mod(float64(i), 2.0) == 1 {
            fmt.Printf(fmtStr + "\n", strings.Repeat("-", lenI0))
        }
    }
    fmt.Printf("\n==================================\n")
}

func GenerateHtml(digitsD d.Digits, quotient d.Digits,
                          interim []d.Digits, rem int) (htmlStr string) {
    // Print the quotient, a space and then remainder = the remainder
    // Print the divide bar
    // Print the denom and a vertical bar then the first row of interim
    // Print the interim rows
    // Get the padding right!

    lenD := digitsD.Len()
    fmtStr := fmt.Sprintf("%%%vv", lenD + 1 + quotient.Len())
    htmlStr = fmt.Sprintf(fmtStr + " remainder=%v<br>", quotient.Sprint(), rem)
    lenI0 := interim[0].Len()
    htmlStr += fmt.Sprintf(fmtStr + "<br>", strings.Repeat("-",lenI0 + 1))
    htmlStr += fmt.Sprintf("%v|%v<br>", digitsD.Sprint(), interim[0].Sprint())
    endStop := lenD + 1 + lenI0
    for i := 1; i < len(interim); i++ {
        if interim[i].Len() < 1 {
            break;
        }
        fmtStr = fmt.Sprintf("%%%vv", endStop)
        htmlStr += fmt.Sprintf(fmtStr + "<br>", interim[i].Sprint())
        if math.Mod(float64(i), 2.0) == 1 {
            htmlStr += fmt.Sprintf(fmtStr + "<br>", strings.Repeat("-", lenI0))
        }
    }
    return htmlStr
}

func SelfTest() (bool) {

    result := true
    fmt.Printf("Self-tests ...\n")
    // Test Init of zero does somethng sensible
    var testDigits d.Digits
    testDigits.Init(0)
    if (testDigits.Len()==1 && testDigits.Digit(0)==0) {
        if DEBUG {testDigits.Print("testDigits")}
    } else {
        result = false
        if DEBUG {fmt.Printf("Zero test failed.\n")}
    }

    // Test a simple case: divides exactly, zero remainder
    digitsD, quotient, interim, _ := LongDiv(6, 36)
    var digitsExp d.Digits
    // Check denominator.
    digitsExp.Init(6)
    result = result && digitsExp.Compare(digitsD)
    // Check the quotient
    digitsExp.InitFixed(2, 6)
    result = result && digitsExp.Compare(quotient)
    // Check the numerator and rows of working
    if len(interim) == 5 {
        digitsExp.Init(36)
        result = result && digitsExp.Compare(interim[0])
        digitsExp.Init(36)
        result = result && digitsExp.Compare(interim[1])
        digitsExp.Init(0)
        result = result && digitsExp.Compare(interim[2])
    } else {
        fmt.Printf("Expected len(interim)=3 but got %v\n", len(interim))
        result = false
    }

    // Now test a case where the denominator divides each digit
    // of the numerator exactly.
    digitsD, quotient, interim, _ = LongDiv(3, 36)
    // Check denominator.
    digitsExp.Init(3)
    result = result && digitsExp.Compare(digitsD)
    // Check the quotient
    digitsExp.InitFixed(2, 12)
    result = result && digitsExp.Compare(quotient)
    // Check the numerator and rows of working
    if len(interim) == 5 {
        digitsExp.Init(36)
        result = result && digitsExp.Compare(interim[0])
        digitsExp.Init(30)
        result = result && digitsExp.Compare(interim[1])
        digitsExp.Init(6)
        result = result && digitsExp.Compare(interim[2])
        digitsExp.Init(6)
        result = result && digitsExp.Compare(interim[3])
        digitsExp.Init(0)
        result = result && digitsExp.Compare(interim[4])
    } else {
        fmt.Printf("Expected len(interim)=3 but got %v\n", len(interim))
        result = false
    }

    // Now test the LongDiv function with some numbers that
    // exercise carrying logic and non-zero remainder
    digitsD, quotient, interim, _ = LongDiv(12, 143)
    // Set up expected result
    digitsExp.Init(12)
    result = result && digitsExp.Compare(digitsD)
    digitsExp.InitFixed(3, 11)
    result = result && digitsExp.Compare(quotient)
    if len(interim) == 7 {
        digitsExp.Init(143)
        result = result && digitsExp.Compare(interim[0])
        digitsExp.Init(120)
        result = result && digitsExp.Compare(interim[1])
        digitsExp.Init(23)
        result = result && digitsExp.Compare(interim[2])
        digitsExp.Init(12)
        result = result && digitsExp.Compare(interim[3])
        digitsExp.Init(11)
        result = result && digitsExp.Compare(interim[4])
    } else {
        fmt.Printf("Expected len(interim)=7 but got %v\n", len(interim))
        result = false
    }

    if result {
        fmt.Printf("... Passed.\n")
    } else {
        fmt.Printf("... FAILED.\n")
    }

    return result
}
