package main

import (
    "fmt"
    "math"
)


func Gcd(ra float64, rb float64) float64 {

    rt := 0.0
    for rb != 0 {
        rt = rb
        rb = math.Mod(ra, rb)
        ra = rt
    }
    return ra
}

func main() {
    gcdAns := Gcd(1001, 1008)
    fmt.Printf("Ans = %9.f\n", gcdAns)
}
