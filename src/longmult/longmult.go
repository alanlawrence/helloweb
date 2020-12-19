package main

import (
    "fmt"
    "math"
)


// Algorithm: https://en.wikipedia.org/wiki/Multiplication_algorithm#Long_multiplication
// TODO: Change code so first arg is the multiplier anum.
// TODO: Stretch goal. Compute the carrys for the addition.
func LongMult(bnum int, anum int) ([]float64, []float64, 
                                   [][]float64, [][]float64) {

    base := 10.0
    //turn ints into a string of int digits
    lenA := 1
    if anum > 0 {
        lenA = int(math.Trunc(math.Log10(float64(anum)) + 1.0))
    }
    lenB := 1
    if bnum > 0 {
        lenB = int(math.Trunc(math.Log10(float64(bnum)) + 1.0))
    }
    digitsA := make([]float64,lenA)
    for di := 0; di < lenA; di++ {
        digitsA[di] = math.Mod(
            math.Trunc(float64(anum)/math.Pow10(lenA-1-di)), base)
    }
    digitsB := make([]float64,lenB)
    for di := 0; di < lenB; di++ {
        digitsB[di] = math.Mod(
            math.Trunc(float64(bnum)/math.Pow10(lenB-1-di)), base)
    }

    //product := array [lenB + 1][lenA + lenB]
    //    A row of interim results for each digit in multiplier bnum
    //    and 1 extra row to hold the Total.
    //    The interim result rows are not necessary to compute the Total
    //    but they are required in long multiplication working.
    product := make([][]float64, lenB + 1)
    for ri := range product {
        product[ri] = make([]float64, lenA + lenB)
    }
    //An array to store the mulitplication carries. 
    //Not necessary for the calculation
    //but required in long multiplication working.
    carrys := make([][]float64, lenB)
    for ri := range carrys {
        carrys[ri] = make([]float64, lenA + lenB)
    }

    carryTotal := 0.0
    carryRow := 0.0
    pi := 0
    // We descend the array from left to right so we read
    // the digits in human readable order.
    for bi := lenB-1; bi >= 0; bi-- {
        carryTotal = 0
        carryRow = 0
        for ai := lenA-1; ai >= 0; ai-- {
            pi = ai + bi + 1
            // Per row interim results (row per digit of mulitplier anum)
            product[bi][pi] = carryRow + digitsA[ai] * digitsB[bi]
            carryRow = math.Trunc(product[bi][pi] / base)
            carrys[bi][pi] = carryRow
            product[bi][pi] = math.Mod(product[bi][pi], base)

            // Total. Note the += to sum the interim results as we go.
            product[lenB][pi] += carryTotal + digitsA[ai] * digitsB[bi]
            carryTotal = math.Trunc(product[lenB][pi] / base)
            product[lenB][pi] = math.Mod(product[lenB][pi], base)
        }
        product[bi][bi] = carryRow
        product[lenB][bi] = carryTotal
    }
    return digitsA, digitsB, product, carrys
}

func PrintWorking(digitsA []float64, digitsB []float64,
                  product [][]float64, carrys [][]float64) {

    lenS := len(digitsA) - len(digitsB)
    for si := 0; si < -1 * lenS; si++ {
        fmt.Printf("  ")
    }
    fmt.Printf("            %v\n", digitsA)
    for si := 0; si < lenS; si++ {
        fmt.Printf("  ")
    }
    fmt.Printf("          x %v\n\n", digitsB)
    lenP := len(product)
    printZeroes := false
    for ri := lenP - 2 ; ri >= 0; ri-- {
        fmt.Printf("          ")
        printZeroes = false
        for _, digitP := range product[ri] {
            if digitP == 0 && !printZeroes {
                fmt.Printf("  ")
            } else {
                fmt.Printf("%v ", digitP)
                printZeroes = true
            }
        }
        fmt.Printf("\n");
        fmt.Printf("carry-> ")
        for _, digitC := range carrys[ri] {
            if digitC == 0 {
                fmt.Printf("  ")
            } else {
                fmt.Printf("%v ", digitC)
            }
        }
        fmt.Printf("\n");
    }
    fmt.Printf("Total     ")
    printZeroes = false
    for _, digitT := range product[lenP-1] {
        if digitT == 0 && !printZeroes {
            fmt.Printf("  ")
        } else {
            fmt.Printf("%v ", digitT)
            printZeroes = true
        }
    }
    fmt.Printf("\n");
}

func main() {
    PrintWorking(LongMult(24, 418))
}
